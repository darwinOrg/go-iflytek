package dgkdxf

import (
	"errors"
	"github.com/darwinOrg/go-common/model"
	"github.com/darwinOrg/go-common/utils"
	"time"
)

type AudioType string

const (
	AudioTypeRaw     AudioType = "raw"
	AudioTypeSpeex   AudioType = "speex"
	AudioTypeOpusOgg AudioType = "opus-ogg"

	timeFormat     = "2006-01-02T15:04:05Z0700"
	apiSuccessCode = "000000"

	defaultBufferSize = 1024 * 16
)

var ApiNoSuccessErr = errors.New("api resp no success")
var ApiGetResultFailTypeErr = errors.New("api get result fail")

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

func getNowTimeString() string {
	return time.Now().Format(timeFormat)
}
