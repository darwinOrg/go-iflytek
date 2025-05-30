package iflytek

import (
	"encoding/json"
	dgctx "github.com/darwinOrg/go-common/context"
	dgerr "github.com/darwinOrg/go-common/enums/error"
	"github.com/darwinOrg/go-common/model"
	"github.com/darwinOrg/go-common/utils"
	dglogger "github.com/darwinOrg/go-logger"
	dgws "github.com/darwinOrg/go-websocket"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"net/url"
	"strings"
	"time"
)

type RoleType int
type AstResultType string
type EngVadMdnType int
type GetBizIdHandler func(ctx *dgctx.DgContext) int64
type SaveAstStartedMetaHandler func(*dgctx.DgContext, string, string) error
type ConsumeAstResultHandler func(*dgctx.DgContext, *AstResult, time.Time) error

const (
	RoleTypeClose RoleType = 0
	RoleTypeOpen  RoleType = 2

	AstResultTypeFinal  AstResultType = "0"
	AstResultTypeMiddle AstResultType = "1"

	EngVadMdnTypeFar  EngVadMdnType = 1
	EngVadMdnTypeNear EngVadMdnType = 2

	ContextIdKey   = "contextId"
	SessionIdKey   = "sessionId"
	CurrentRoleKey = "currentRole"

	ExceedUploadSpeedLimitCode = "100001"
	UnknownErrorCode           = "999999"
)

type AstParamConfig struct {
	Lang           string        `json:"lang"`
	Codec          string        `json:"codec"`
	AudioEncode    string        `json:"audioEncode"`
	Samplerate     string        `json:"samplerate"`
	RoleType       RoleType      `json:"roleType"`
	ContextId      string        `json:"contextId"`
	FeatureIds     string        `json:"featureIds"`
	HotWordId      string        `json:"hotWordId"`
	SourceInfo     string        `json:"sourceInfo"`
	FilePath       string        `json:"filePath"`
	ResultFilePath string        `json:"resultFilePath"`
	EngVadMdn      EngVadMdnType `json:"engVadMdn" remark:"vad 远近场切换。不传此参数或传值1代表远场。若要使用近场，传值2"`
}

type AstReadMessageRequest struct {
	ForwardMark               string
	BizKey                    string
	BizIdHandler              GetBizIdHandler
	SaveAstStartedMetaHandler SaveAstStartedMetaHandler
	ConsumeAstResultHandler   ConsumeAstResultHandler
}

type AstResult struct {
	SegID int64 `json:"seg_id"`
	Cn    struct {
		St struct {
			Bg string `json:"bg"`
			Ed string `json:"ed"`
			Rt []struct {
				Ws []struct {
					Cw []struct {
						W  string `json:"w"`
						Wp string `json:"wp"`
						Rl string `json:"rl"`
					} `json:"cw"`
					Wb int64 `json:"wb"`
					We int64 `json:"we"`
				} `json:"ws"`
			} `json:"rt"`
			Type AstResultType `json:"type"`
		} `json:"st"`
	} `json:"cn"`
	Ls bool `json:"ls"`
}

func (ar *AstResult) HasFinalWords() bool {
	return ar.Cn.St.Type == AstResultTypeFinal && len(ar.Cn.St.Rt) > 0
}

func (ar *AstResult) CombineFinalWords(ctx *dgctx.DgContext, roleType RoleType) string {
	var finalWords string

	if ar.HasFinalWords() {
		for _, rt := range ar.Cn.St.Rt {
			if len(rt.Ws) > 0 {
				for _, ws := range rt.Ws {
					if len(ws.Cw) > 0 {
						for _, cw := range ws.Cw {
							finalWords = finalWords + cw.W

							if roleType == RoleTypeOpen {
								SetCurrentRole(ctx, cw.Rl)
							}
						}
					}
				}
			}
		}
	}

	return finalWords
}

func (c *Client) AstConnect(ctx *dgctx.DgContext, config *AstParamConfig) (*websocket.Conn, error) {
	if int(config.EngVadMdn) == 0 {
		config.EngVadMdn = EngVadMdnTypeFar
	}
	uri := c.BuildAstUri(ctx, config)
	dglogger.Infof(ctx, "ast config: %s, uri: %s", utils.MustConvertBeanToJsonString(config), uri)
	cn, _, err := websocket.DefaultDialer.Dial(uri, nil)
	if err != nil {
		return nil, err
	}

	return cn, nil
}

func (c *Client) BuildAstUri(ctx *dgctx.DgContext, config *AstParamConfig) string {
	parts := []string{"v1.0", c.Config.AppId, c.Config.AccessKeyId, getDateTimeString(), uuid.NewString()}
	partsStr := strings.Join(parts, ",")
	baseString := url.QueryEscape(partsStr)
	signature := utils.Sha1Base64Encode(c.Config.AccessKeySecret, baseString)
	parts = append(parts, signature)
	authString := strings.Join(parts, ",")

	params := []*model.KeyValuePair[string, any]{
		{
			Key:   "lang",
			Value: config.Lang,
		},
		{
			Key:   "codec",
			Value: config.Codec,
		},
		{
			Key:   "samplerate",
			Value: config.Samplerate,
		},
		{
			Key:   "hotWordId",
			Value: config.HotWordId,
		},
		{
			Key:   "sourceInfo",
			Value: config.SourceInfo,
		},
		{
			Key:   "audioEncode",
			Value: config.AudioEncode,
		},
		{
			Key:   "roleType",
			Value: config.RoleType,
		},
		{
			Key:   "featureIds",
			Value: config.FeatureIds,
		},
		{
			Key:   "eng_vad_mdn",
			Value: config.EngVadMdn,
		},
		{
			Key:   "authString",
			Value: authString,
		},
		{
			Key:   "trackId",
			Value: ctx.TraceId,
		},
	}

	return c.Config.Host + "/ast?" + utils.FormUrlEncodedParams(params)
}

func AstWriteStarted(ctx *dgctx.DgContext, cn *websocket.Conn) error {
	dglogger.Infof(ctx, "send ast started message")
	if cn == nil {
		dglogger.Error(ctx, "websocket conn is nil")
		return dgerr.SYSTEM_ERROR
	}
	return cn.WriteMessage(websocket.TextMessage, []byte("{\"action\":\"started\"}"))
}

func AstWriteEnd(ctx *dgctx.DgContext, cn *websocket.Conn) error {
	dglogger.Infof(ctx, "send ast end message")
	if cn == nil {
		dglogger.Error(ctx, "websocket conn is nil")
		return dgerr.SYSTEM_ERROR
	}
	return cn.WriteMessage(websocket.TextMessage, []byte("{\"end\":true}"))
}

func IsAstEndMessage(_ *dgctx.DgContext, mt int, data []byte) bool {
	if mt == websocket.CloseMessage || mt == -1 {
		return true
	}

	if mt == websocket.TextMessage && len(data) > 0 {
		var mp map[string]any
		err := json.Unmarshal(data, &mp)
		if err != nil {
			return false
		}

		end, ok := mp["end"].(bool)
		if ok && end {
			return true
		}
	}

	return false
}

func SetContextId(ctx *dgctx.DgContext, connMark string, contextId string) bool {
	if contextId == "" {
		return false
	}

	oriContextId := GetContextId(ctx, connMark)
	if oriContextId != "" && oriContextId == contextId {
		return false
	}

	ctx.SetExtraKeyValue(ContextIdKey+connMark, contextId)
	return true
}

func GetContextId(ctx *dgctx.DgContext, connMark string) string {
	ctxId := ctx.GetExtraValue(ContextIdKey + connMark)
	if ctxId == nil {
		return ""
	}

	return ctxId.(string)
}

func SetSessionId(ctx *dgctx.DgContext, connMark string, sessionId string) {
	ctx.SetExtraKeyValue(SessionIdKey+connMark, sessionId)
}

func GetSessionId(ctx *dgctx.DgContext, connMark string) string {
	sessionId := ctx.GetExtraValue(SessionIdKey + connMark)
	if sessionId == nil {
		return ""
	}

	return sessionId.(string)
}

func SetCurrentRole(ctx *dgctx.DgContext, currentRole string) bool {
	if currentRole == "" || currentRole == "0" {
		return false
	}

	oriCurrentRole := GetCurrentRole(ctx)
	if oriCurrentRole != "" && oriCurrentRole == currentRole {
		return false
	}

	ctx.SetExtraKeyValue(CurrentRoleKey, currentRole)
	return true
}

func GetCurrentRole(ctx *dgctx.DgContext) string {
	currentRole := ctx.GetExtraValue(CurrentRoleKey)
	if currentRole == nil {
		return ""
	}

	return currentRole.(string)
}

func AstReadMessage(ctx *dgctx.DgContext, req *AstReadMessageRequest) {
	dgws.IncrWaitGroup(ctx)
	defer dgws.DoneWaitGroup(ctx)
	bizId := req.BizIdHandler(ctx)
	forwardMark := req.ForwardMark
	bizKey := req.BizKey

	for {
		if dgws.IsWsEnded(ctx) {
			dglogger.Infof(ctx, "[%s: %d, forwardMark: %s] websocket already ended", bizKey, bizId, forwardMark)
			return
		}

		if dgws.IsForwardWsEnded(ctx, forwardMark) {
			time.Sleep(time.Second)
			continue
		}

		forwardConn := dgws.GetForwardConn(ctx, forwardMark)
		if forwardConn == nil {
			dglogger.Debugf(ctx, "[%s: %d, forwardMark: %s] forward conn is nil", bizKey, bizId, forwardMark)
			time.Sleep(time.Second)
			continue
		}
		mt, data, err := forwardConn.ReadMessage()
		if mt == websocket.CloseMessage || mt == -1 {
			dglogger.Infof(ctx, "[%s: %d, forwardMark: %s] received iflytek ast close message, error: %v", bizKey, bizId, forwardMark, err)
			dgws.SetForwardWsEnded(ctx, forwardMark)
			continue
		}

		if mt == websocket.TextMessage {
			dglogger.Debugf(ctx, "[%s: %d, forwardMark: %s] receive iflytek ast message: %s", bizKey, bizId, forwardMark, string(data))
			var mp map[string]any
			err := json.Unmarshal(data, &mp)
			if err != nil {
				dglogger.Errorf(ctx, "[%s: %d, forwardMark: %s] unmarshal message[%s] error: %v", bizKey, bizId, forwardMark, string(data), err)
				continue
			}

			action := mp["action"]
			if action == "started" {
				dglogger.Infof(ctx, "[%s: %d, forwardMark: %s] received iflytek ast started message", bizKey, bizId, forwardMark)
				if req.SaveAstStartedMetaHandler != nil {
					contextId := mp[ContextIdKey].(string)
					sessionId := mp[SessionIdKey].(string)
					err := req.SaveAstStartedMetaHandler(ctx, contextId, sessionId)
					if err != nil {
						dglogger.Errorf(ctx, "[%s: %d, forwardMark: %s] save ast started meta[contextId: %s, sessionId: %s] error: %v", bizKey, bizId, forwardMark, contextId, sessionId, err)
					}
				}

				continue
			}

			code := mp["code"]
			if code == ExceedUploadSpeedLimitCode {
				dglogger.Errorf(ctx, "[%s: %d, forwardMark: %s] iflytek ast exceed upload speed limit", bizKey, bizId, forwardMark)
				continue
			}

			astResult, err := utils.ConvertJsonBytesToBean[AstResult](data)
			if err != nil {
				dglogger.Errorf(ctx, "[%s: %d, forwardMark: %s] unmarshal message[%s] error: %v", bizKey, bizId, forwardMark, string(data), err)
				continue
			}

			if astResult != nil && req.ConsumeAstResultHandler != nil {
				err := req.ConsumeAstResultHandler(ctx, astResult, time.Now())
				if err != nil {
					dglogger.Errorf(ctx, "[%s: %d, forwardMark: %s] consume ast message[%s] error: %v", bizKey, bizId, forwardMark, string(data), err)
				}
			}

			continue
		}

		if err != nil {
			dglogger.Errorf(ctx, "[%s: %d, forwardMark: %s] read ast message error: %v", bizKey, bizId, forwardMark, err)
		}
	}
}
