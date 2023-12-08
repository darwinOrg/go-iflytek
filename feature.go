package dgkdxf

import (
	"errors"
	dgctx "github.com/darwinOrg/go-common/context"
	"github.com/darwinOrg/go-common/model"
	"github.com/darwinOrg/go-common/utils"
	dghttp "github.com/darwinOrg/go-httpclient"
	"github.com/google/uuid"
	"strings"
)

const featureFailStatus = 2

type FeatureResult[T any] struct {
	Code string `json:"code"`
	Desc string `json:"desc"`
	Data *T     `json:"data"`
	Sid  string `json:"sid"`
}

func (rt *FeatureResult[T]) isSuccess() bool {
	return rt.Code == apiSuccessCode
}

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
	params, header := c.buildFeatureParamsAndHeader(ctx)
	url := c.Config.Host + "/res/feature/v1/register?" + utils.FormUrlEncodedParams(params)
	dghttp.SetHttpClient(ctx, dghttp.Client11)
	rt, err := dghttp.DoPostJsonToStruct[FeatureResult[string]](ctx, url, req, header)
	if err != nil {
		return "", err
	}

	if !rt.isSuccess() {
		return "", errors.New(rt.Desc)
	}

	resp, err := utils.ConvertJsonStringToBean[RegisterFeatureResponse](*rt.Data)
	if err != nil {
		return "", err
	}

	if resp.Status == featureFailStatus || resp.FeatureId == "" {
		return "", errors.New("注册失败")
	}

	return resp.FeatureId, nil
}

func (c *Client) UpdateFeature(ctx *dgctx.DgContext, req *UpdateFeatureRequest) error {
	params, header := c.buildFeatureParamsAndHeader(ctx)
	url := c.Config.Host + "/res/feature/v1/update?" + utils.FormUrlEncodedParams(params)
	dghttp.SetHttpClient(ctx, dghttp.Client11)
	rt, err := dghttp.DoPostJsonToStruct[FeatureResult[string]](ctx, url, req, header)
	if err != nil {
		return err
	}

	if !rt.isSuccess() {
		return errors.New(rt.Desc)
	}

	resp, err := utils.ConvertJsonStringToBean[UpdateFeatureResponse](*rt.Data)
	if err != nil {
		return err
	}

	if resp.Status == featureFailStatus {
		return errors.New("更新失败")
	}

	return nil
}

func (c *Client) DeleteFeature(ctx *dgctx.DgContext, featureIds []string) []string {
	req := map[string]any{"feature_ids": featureIds}
	params, header := c.buildFeatureParamsAndHeader(ctx)
	url := c.Config.Host + "/res/feature/v1/update?" + utils.FormUrlEncodedParams(params)
	dghttp.SetHttpClient(ctx, dghttp.Client11)
	rt, err := dghttp.DoPostJsonToStruct[FeatureResult[string]](ctx, url, req, header)
	if err != nil {
		return featureIds
	}

	if !rt.isSuccess() {
		return featureIds
	}

	resp, err := utils.ConvertJsonStringToBean[DeleteFeatureResponse](*rt.Data)
	if err != nil {
		return featureIds
	}

	if resp.DelFailIds == "" {
		return []string{}
	}

	return strings.Split(resp.DelFailIds, ";")
}

func (c *Client) buildFeatureParamsAndHeader(ctx *dgctx.DgContext) ([]*model.KeyValuePair[string, any], map[string]string) {
	params := []*model.KeyValuePair[string, any]{
		{
			Key:   "accessKeyId",
			Value: c.Config.AccessKeyId,
		},
		{
			Key:   "dateTime",
			Value: getNowTimeString(),
		},
		{
			Key:   "signatureRandom",
			Value: uuid.NewString(),
		},
	}

	header := map[string]string{
		"signature": c.GenerateSignature(params),
		"x-traceid": ctx.TraceId,
	}

	return params, header
}
