package dgkdxf

import (
	"encoding/json"
	dgcoll "github.com/darwinOrg/go-common/collection"
	dgctx "github.com/darwinOrg/go-common/context"
	"github.com/darwinOrg/go-common/model"
	"github.com/darwinOrg/go-common/utils"
	dglogger "github.com/darwinOrg/go-logger"
	dgws "github.com/darwinOrg/go-websocket"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type RoleType int
type AstResultType string
type GetBizIdFunc func(ctx *dgctx.DgContext) int64
type SaveAstStartedMetaFunc func(*dgctx.DgContext, string, string) error
type ConsumeAstResultFunc func(*dgctx.DgContext, *AstResult, time.Time) error

const (
	RoleTypeClose RoleType = 0
	RoleTypeOpen  RoleType = 2

	AstResultTypeFinal  AstResultType = "0"
	AstResultTypeMiddle AstResultType = "1"

	ContextIdKey   = "contextId"
	SessionIdKey   = "sessionId"
	CurrentRoleKey = "currentRole"

	ExceedUploadSpeedLimitCode = "100001"
	UnknownErrorCode           = "999999"
)

var punctuations = []string{",", ".", "?", "!", ";", ":", "'", "\"", "(", ")", "{", "}", "[", "]", "<", ">", "@", "#", "$", "%", "^", "&", "*", "+", "=", "-", "_", "|", "~", "，", "。", "？", "！", "；", "：", "“", "”", "‘", "’", "《", "》", "（", "）", "【", "】"}

type AstParamConfig struct {
	Lang           string   `json:"lang"`
	Codec          string   `json:"codec"`
	AudioEncode    string   `json:"audioEncode"`
	Samplerate     string   `json:"samplerate"`
	RoleType       RoleType `json:"roleType"`
	ContextId      string   `json:"contextId"`
	FeatureIds     string   `json:"featureIds"`
	HotWordId      string   `json:"hotWordId"`
	SourceInfo     string   `json:"sourceInfo"`
	FilePath       string   `json:"filePath"`
	ResultFilePath string   `json:"resultFilePath"`
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

func (ar *AstResult) CombineFinalWords(ctx *dgctx.DgContext) string {
	var finalWords string

	if ar.HasFinalWords() {
		for _, rt := range ar.Cn.St.Rt {
			if len(rt.Ws) > 0 {
				for _, ws := range rt.Ws {
					if len(ws.Cw) > 0 {
						for _, cw := range ws.Cw {
							finalWords = finalWords + cw.W
							SetCurrentRole(ctx, cw.Rl)
						}
					}
				}
			}
		}
	}

	return deleteStartPunctuation(finalWords)
}

func (c *Client) AstConnect(ctx *dgctx.DgContext, config *AstParamConfig) (*websocket.Conn, error) {
	uri := c.BuildAstUri(ctx, config)
	dglogger.Infof(ctx, "ast config: %s, uri: %s", utils.MustConvertBeanToJsonString(config), uri)
	cn, _, err := websocket.DefaultDialer.Dial(uri, nil)
	if err != nil {
		return nil, err
	}

	return cn, nil
}

func (c *Client) BuildAstUri(ctx *dgctx.DgContext, config *AstParamConfig) string {
	parts := []string{"v1.0", c.Config.AppId, c.Config.AccessKeyId, getNowTimeString(), uuid.NewString()}
	partsStr := strings.Join(parts, ",")
	baseString := url.QueryEscape(partsStr)
	signature := utils.Sha1Base64Encode(c.Config.AccessKeySecret, baseString)
	parts = append(parts, signature)
	authString := strings.Join(parts, ",")

	params := []*model.KeyValuePair[string, string]{
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
			Value: strconv.Itoa(int(config.RoleType)),
		},
		{
			Key:   "featureIds",
			Value: config.FeatureIds,
		},

		{
			Key:   "authString",
			Value: url.QueryEscape(authString),
		},
		{
			Key:   "trackId",
			Value: ctx.TraceId,
		},
	}

	paramsArr := dgcoll.MapToList(params, func(p *model.KeyValuePair[string, string]) string {
		return p.Key + "=" + p.Value
	})
	paramsStr := strings.Join(paramsArr, "&")
	return c.Config.Host + "/ast?" + paramsStr
}

func AstWriteStarted(ctx *dgctx.DgContext, cn *websocket.Conn) error {
	dglogger.Infof(ctx, "send ast started message")
	return cn.WriteMessage(websocket.TextMessage, []byte("{\"action\":\"started\"}"))
}

func AstWriteEnd(ctx *dgctx.DgContext, cn *websocket.Conn) error {
	dglogger.Infof(ctx, "send ast end message")
	return cn.WriteMessage(websocket.TextMessage, []byte("{\"end\":true}"))
}

func IsAstEndMessage(mt int, data []byte) bool {
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

func SetContextId(ctx *dgctx.DgContext, contextId string) bool {
	if contextId == "" {
		return false
	}

	oriContextId := GetContextId(ctx)
	if oriContextId != "" && oriContextId == contextId {
		return false
	}

	ctx.SetExtraKeyValue(ContextIdKey, contextId)
	return true
}

func GetContextId(ctx *dgctx.DgContext) string {
	ctxId := ctx.GetExtraValue(ContextIdKey)
	if ctxId == nil {
		return ""
	}

	return ctxId.(string)
}

func SetSessionId(ctx *dgctx.DgContext, sessionId string) {
	ctx.SetExtraKeyValue(SessionIdKey, sessionId)
}

func GetSessionId(ctx *dgctx.DgContext) string {
	sessionId := ctx.GetExtraValue(SessionIdKey)
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

func AstReadMessage(ctx *dgctx.DgContext, conn *websocket.Conn, forwardConn *websocket.Conn, bizKey string, getBizIdFunc GetBizIdFunc, saveAstStartedMetaFunc SaveAstStartedMetaFunc, consumeAstResultFunc ConsumeAstResultFunc) {
	bizId := getBizIdFunc(ctx)

	for {
		if dgws.IsWsEnded(ctx) {
			dglogger.Infof(ctx, "[%s: %d] websocket already ended", bizKey, bizId)
			return
		}

		mt, data, err := forwardConn.ReadMessage()
		if mt == websocket.CloseMessage || mt == -1 {
			dglogger.Infof(ctx, "[%s: %d] received iflytek ast close message, error: %v", bizKey, bizId, err)
			dgws.SetWsEnded(ctx)
			conn.WriteMessage(websocket.CloseMessage, data)
			return
		}
		conn.WriteMessage(mt, data)

		if mt == websocket.TextMessage {
			dglogger.Infof(ctx, "[%s: %d] receive iflytek ast message: %s", bizKey, bizId, string(data))
			var mp map[string]any
			err := json.Unmarshal(data, &mp)
			if err != nil {
				dglogger.Errorf(ctx, "[%s: %d] unmarshal message[%s] error: %v", bizKey, bizId, string(data), err)
				continue
			}

			action := mp["action"]
			if action == "started" {
				dglogger.Infof(ctx, "[%s: %d] received iflytek ast started message", bizKey, bizId)
				if saveAstStartedMetaFunc != nil {
					contextId := mp[ContextIdKey].(string)
					sessionId := mp[SessionIdKey].(string)
					go func() {
						err := saveAstStartedMetaFunc(ctx, contextId, sessionId)
						if err != nil {
							dglogger.Errorf(ctx, "[%s: %d] save ast started meta[contextId: %s, sessionId: %s] error: %v", bizKey, bizId, contextId, sessionId, err)
						}
					}()
				}

				continue
			}

			code := mp["code"]
			if code == ExceedUploadSpeedLimitCode {
				dglogger.Errorf(ctx, "[%s: %d] iflytek ast exceed upload speed limit", bizKey, bizId)
				continue
			}

			astResult, err := utils.ConvertJsonBytesToBean[AstResult](data)
			if err != nil {
				dglogger.Errorf(ctx, "[%s: %d] unmarshal message[%s] error: %v", bizKey, bizId, string(data), err)
				continue
			}

			if astResult != nil && consumeAstResultFunc != nil {
				if astResult.HasFinalWords() {
					go func() {
						err := consumeAstResultFunc(ctx, astResult, time.Now())
						if err != nil {
							dglogger.Errorf(ctx, "[%s: %d] consume ast message[%s] error: %v", bizKey, bizId, string(data), err)
						}
					}()
				}
			}

			continue
		}

		if err != nil {
			dglogger.Errorf(ctx, "[%s: %d] read ast message error: %v", bizKey, bizId, err)
		}
	}
}

func deleteStartPunctuation(str string) string {
	if str == "" {
		return ""
	}

	for _, p := range punctuations {
		_, after, found := strings.Cut(str, p)
		if found {
			return after
		}
	}

	return str
}
