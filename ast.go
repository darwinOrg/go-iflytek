package dgkdxf

import (
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

	actionStarted = "started"
	actionEnd     = "end"
)

type AstParamConfig struct {
	Lang           string   `json:"lang"`
	Codec          string   `json:"codec"`
	Samplerate     string   `json:"samplerate"`
	AudioEncode    string   `json:"audioEncode"`
	HotWordId      string   `json:"hotWordId"`
	SourceInfo     string   `json:"sourceInfo"`
	RoleType       RoleType `json:"roleType"`
	FeatureIds     []string `json:"featureIds"`
	FilePath       string   `json:"filePath"`
	ResultFilePath string   `json:"resultFilePath"`
}

func (c *Client) AstConnect(ctx *dgctx.DgContext, config *AstParamConfig) (*websocket.Conn, error) {
	uri := c.buildAstUri(ctx, config)
	dglogger.Infof(ctx, "ast config: %s, uri: %s", utils.MustConvertBeanToJsonString(config), uri)
	cn, _, err := websocket.DefaultDialer.Dial(uri, nil)
	if err != nil {
		return nil, err
	}

	return cn, nil
}

func (c *Client) AstReadMessage(ctx *dgctx.DgContext, cn *websocket.Conn) error {
	for {
		mt, message, err := cn.ReadMessage()

		if mt == websocket.CloseMessage || mt == -1 {
			dglogger.Infof(ctx, "received iflytek close message, error: %v", err)
			return nil
		}

		if mt == websocket.PongMessage {
			dglogger.Info(ctx, "received iflytek pong message")
			continue
		}

		if mt == websocket.TextMessage {

			continue
		}

		if err != nil {
			return err
		}
	}
}

func (c *Client) buildAstUri(ctx *dgctx.DgContext, config *AstParamConfig) string {
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
			Value: config.Lang,
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
