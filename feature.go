package dgkdxf

import (
	"errors"
	dgctx "github.com/darwinOrg/go-common/context"
	dghttp "github.com/darwinOrg/go-httpclient"
	"github.com/google/uuid"
	"strings"
	"time"
)

type AudioType string

const (
	AudioTypeRaw     AudioType = "raw"
	AudioTypeSpeex   AudioType = "speex"
	AudioTypeOpusOgg AudioType = "opus-ogg"
	failStatus                 = 2
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

type DeleteFeatureRequest struct {
	FeatureIds []string `json:"feature_ids" binding:"required,minLength=1"`
}

type DeleteFeatureResponse struct {
	DelFailIds string `json:"del_fail_ids"`
}

func (c *Client) RegisterFeature(ctx *dgctx.DgContext, req *RegisterFeatureRequest) (string, error) {
	params, header := c.buildParamsAndHeader(ctx)
	url := c.Config.Host + "/res/feature/v1/register?" + formUrlEncodedParams(params)
	rt, err := dghttp.DoPostJsonToStruct[IflytekResult[RegisterFeatureResponse]](ctx, url, req, header)
	if err != nil {
		return "", err
	}

	if !rt.isSuccess() {
		return "", errors.New(rt.Desc)
	}

	if rt.Data.Status == failStatus || rt.Data.FeatureId == "" {
		return "", errors.New("注册失败")
	}

	return rt.Data.FeatureId, nil
}

func (c *Client) UpdateFeature(ctx *dgctx.DgContext, req *UpdateFeatureRequest) error {
	params, header := c.buildParamsAndHeader(ctx)
	url := c.Config.Host + "/res/feature/v1/update?" + formUrlEncodedParams(params)
	rt, err := dghttp.DoPostJsonToStruct[IflytekResult[UpdateFeatureResponse]](ctx, url, req, header)
	if err != nil {
		return err
	}

	if !rt.isSuccess() {
		return errors.New(rt.Desc)
	}

	if rt.Data.Status == failStatus {
		return errors.New("更新失败")
	}

	return nil
}

func (c *Client) DeleteFeature(ctx *dgctx.DgContext, featureIds []string) []string {
	req := map[string]any{"feature_ids": featureIds}
	params, header := c.buildParamsAndHeader(ctx)
	url := c.Config.Host + "/res/feature/v1/update?" + formUrlEncodedParams(params)
	rt, err := dghttp.DoPostJsonToStruct[IflytekResult[DeleteFeatureResponse]](ctx, url, req, header)
	if err != nil {
		return featureIds
	}

	if !rt.isSuccess() {
		return featureIds
	}

	if rt.Data.DelFailIds == "" {
		return []string{}
	}

	return strings.Split(rt.Data.DelFailIds, ";")
}

func (c *Client) buildParamsAndHeader(ctx *dgctx.DgContext) (map[string]any, map[string]string) {
	params := map[string]any{
		"accessKeyId":     c.Config.AccessKeyId,
		"dateTime":        time.Now().Format(timeFormat),
		"signatureRandom": uuid.NewString(),
	}

	header := map[string]string{
		"signature": c.GenerateSignature(params),
		"x-traceid": ctx.TraceId,
	}

	return params, header
}
