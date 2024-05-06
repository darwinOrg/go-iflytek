package dgkdxf

import (
	"errors"
	dgcoll "github.com/darwinOrg/go-common/collection"
	"github.com/darwinOrg/go-common/model"
	"github.com/darwinOrg/go-common/utils"
	"time"
)

type AudioType string

const (
	AudioTypeRaw     AudioType = "raw"
	AudioTypeSpeex   AudioType = "speex"
	AudioTypeOpusOgg AudioType = "opus-ogg"
	AudioTypeOpusWb  AudioType = "opus-wb"

	dateTimeFormat  = "2006-01-02T15:04:05Z0700"
	timestampFormat = "2006-01-02T15:04:05Z"
	apiSuccessCode  = "000000"

	defaultBufferSize = 1024 * 16
)

var (
	ApiNoSuccessErr         = errors.New("api resp no success")
	ApiGetResultFailTypeErr = errors.New("api get result fail")
)

type KdxfResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
	RequestID string `json:"requestId"`
}

type ClientConfig struct {
	AppId           string `json:"appId"`
	Host            string `json:"host"`
	AccessKeyId     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
}

type Client struct {
	Config *ClientConfig
}

func NewClient(config *ClientConfig) *Client {
	return &Client{Config: config}
}

func (c *Client) GenerateSignature(params []*model.KeyValuePair[string, any]) string {
	baseString := utils.FormUrlEncodedParams(params)

	return utils.Sha1Base64Encode(c.Config.AccessKeySecret, baseString)
}

func (c *Client) GenerateSignatureWithUrlPrefix(urlPrefix string, params []*model.KeyValuePair[string, any]) string {
	baseString := utils.FormUrlEncodedParams(params)

	return utils.Sha1Base64Encode(c.Config.AccessKeySecret, urlPrefix+baseString)
}

func (c *Client) buildCommonParams() []*model.KeyValuePair[string, any] {
	params := []*model.KeyValuePair[string, any]{
		{
			Key:   "Timestamp",
			Value: getTimestampString(),
		},
		{
			Key:   "AccessKeyId",
			Value: c.Config.AccessKeyId,
		},
		{
			Key:   "Expires",
			Value: 86400,
		},
	}

	sortParams(params)
	return params
}

func getDateTimeString() string {
	return time.Now().Format(dateTimeFormat)
}

func getTimestampString() string {
	return time.Now().Format(timestampFormat)
}

func sortParams(params []*model.KeyValuePair[string, any]) {
	dgcoll.SortAsc(params, func(p *model.KeyValuePair[string, any]) string { return p.Key })
}
