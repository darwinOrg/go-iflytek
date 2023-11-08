package dgkdxf

import (
	"errors"
	dgctx "github.com/darwinOrg/go-common/context"
	dghttp "github.com/darwinOrg/go-httpclient"
	"github.com/google/uuid"
	"time"
)

const (
	registerFeaturePath = "/res/feature/v1/register"
	updateFeaturePath   = "/res/feature/v1/update"
	deleteFeaturePath   = "/res/feature/v1/delete"
)

type AudioType string

const (
	AudioTypeRaw     = "raw"
	AudioTypeSpeex   = "speex"
	AudioTypeOpusOgg = "opus-ogg"
)

type RegisterFeatureRequest struct {
	AudioData string    `json:"audio_data" binding:"required,minLength=1"`
	AudioType AudioType `json:"audio_type" binding:"required,mustIn=raw#speex#opus-ogg"`
	Uid       string    `json:"uid"`
}

type RegisterFeatureResponse struct {
	FeatureId string `json:"feature_id"`
	Status    int    `json:"status"`
}

type UpdateFeatureRequest struct {
	FeatureId string    `json:"feature_id" binding:"required"`
	AudioData string    `json:"audio_data" binding:"required,minLength=1"`
	AudioType AudioType `json:"audio_type" binding:"required,mustIn=raw#speex#opus-ogg"`
}

type UpdateFeatureResponse struct {
	Status int `json:"status"`
}

func (c *Client) RegisterFeature(ctx *dgctx.DgContext, req *RegisterFeatureRequest) (string, error) {
	params := map[string]any{
		"accessKeyId":     c.Config.AccessKeyId,
		"dateTime":        time.Now().Format(timeFormat),
		"signatureRandom": uuid.NewString(),
	}

	header := map[string]string{
		"signature": c.GenerateSignature(params),
		"x-traceid": ctx.TraceId,
	}

	url := c.Config.Host + "?" + formUrlEncodedParams(params)
	rt, err := dghttp.DoPostJsonToStruct[IflytekResult[RegisterFeatureResponse]](ctx, url, req, header)
	if err != nil {
		return "", err
	}

	if !rt.isSuccess() {
		return "", errors.New(rt.Desc)
	}

	if rt.Data.Status == 2 || rt.Data.FeatureId == "" {
		return "", errors.New("注册失败")
	}

	return rt.Data.FeatureId, nil
}

func (c *Client) UpdateFeature(req *UpdateFeatureRequest) error {
	return nil
}

func (c *Client) DeleteFeature(featureIds []string) []string {
	return nil
}
