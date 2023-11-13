package dgkdxf

import (
	"encoding/json"
	dgcoll "github.com/darwinOrg/go-common/collection"
	dgctx "github.com/darwinOrg/go-common/context"
	"github.com/darwinOrg/go-common/model"
	"github.com/darwinOrg/go-common/utils"
	dglogger "github.com/darwinOrg/go-logger"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"net/url"
	"strconv"
	"strings"
)

type RoleType int

const (
	RoleTypeClose RoleType = 0
	RoleTypeOpen  RoleType = 2
)

type AstParamConfig struct {
	Lang           string   `json:"lang"`
	Codec          string   `json:"codec"`
	AudioEncode    string   `json:"audioEncode"`
	Samplerate     string   `json:"samplerate"`
	RoleType       RoleType `json:"roleType"`
	ContextId      string   `json:"contextId"`
	FeatureIds     []string `json:"featureIds"`
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
						Rl int64  `json:"rl"`
						W  string `json:"w"`
						Wp string `json:"wp"`
					} `json:"cw"`
					Wb int64 `json:"wb"`
					We int64 `json:"we"`
				} `json:"ws"`
			} `json:"rt"`
			Type string `json:"type"`
		} `json:"st"`
	} `json:"cn"`
	Ls bool `json:"ls"`
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
			Value: strings.Join(config.FeatureIds, ","),
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

func AstReadMessage(ctx *dgctx.DgContext, cn *websocket.Conn, consumeFunc func(*dgctx.DgContext, *AstResult) error) error {
	for {
		mt, data, err := cn.ReadMessage()

		if mt == websocket.CloseMessage || mt == -1 {
			dglogger.Infof(ctx, "[userId: %d] received iflytek ast close message, error: %v", ctx.UserId, err)
			return nil
		}

		if mt == websocket.TextMessage {
			var mp map[string]any
			err := json.Unmarshal(data, &mp)
			if err != nil {
				dglogger.Errorf(ctx, "[userId: %d] unmarshal message[%s] error: %v", ctx.UserId, string(data), err)
				continue
			}

			action := mp["action"]
			if action == "started" {
				dglogger.Infof(ctx, "[userId: %d] received iflytek ast started message", ctx.UserId)
				continue
			}

			code := mp["code"]
			if code == "100001" {
				dglogger.Errorf(ctx, "[userId: %d] iflytek ast exceed upload speed limit", ctx.UserId)
				continue
			}

			astResult, err := utils.ConvertJsonBytesToBean[AstResult](data)
			if err != nil {
				dglogger.Errorf(ctx, "[userId: %d] unmarshal message[%s] error: %v", ctx.UserId, string(data), err)
				continue
			}

			if astResult != nil {
				err := consumeFunc(ctx, astResult)
				if err != nil {
					dglogger.Errorf(ctx, "[userId: %d] handle message[%s] error: %v", ctx.UserId, string(data), err)
					continue
				}
			}

			continue
		}

		if err != nil {
			return err
		}
	}
}

func AstWriteStarted(ctx *dgctx.DgContext, cn *websocket.Conn) error {
	dglogger.Infof(ctx, "send started message")
	return cn.WriteMessage(websocket.TextMessage, []byte("{\"action\":\"started\"}"))
}

func AstWriteEnd(ctx *dgctx.DgContext, cn *websocket.Conn) error {
	dglogger.Infof(ctx, "send end message")
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
