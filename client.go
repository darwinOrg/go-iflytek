package dgkdxf

const (
	timeFormat     = "2006-01-02T15:04:05Z0700"
	apiSuccessCode = "000000"
)

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

type IflytekResult[T any] struct {
	Code string `json:"code"`
	Desc string `json:"desc"`
	Data *T     `json:"data"`
	Sid  string `json:"sid"`
}

func (rt *IflytekResult[T]) isSuccess() bool {
	return rt.Code == apiSuccessCode
}
