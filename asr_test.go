package dgkdxf_test

import (
	dgctx "github.com/darwinOrg/go-common/context"
	"github.com/darwinOrg/go-common/utils"
	dgkdxf "github.com/darwinOrg/go-iflytek"
	dglogger "github.com/darwinOrg/go-logger"
	"os"
	"testing"
)

func TestAsrUpload(t *testing.T) {
	host := "https://api.iflyrec.com"
	appId := os.Getenv("appId")
	accessKeyId := os.Getenv("accessKeyId")
	accessKeySecret := os.Getenv("accessKeySecret")
	ctx := &dgctx.DgContext{TraceId: "123"}

	client := dgkdxf.NewClient(&dgkdxf.ClientConfig{
		AppId:           appId,
		Host:            host,
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
	})

	fileBytes, _ := os.ReadFile("test.opus")
	rt, err := client.AsrUpload(ctx, "test.opus", 45370, int64(len(fileBytes)), "")
	if err != nil {
		panic(err)
	}
	dglogger.Infof(ctx, "rt: %s", utils.MustConvertBeanToJsonString(rt))
}

func TestAsrGetResult(t *testing.T) {
	host := "https://api.iflyrec.com"
	appId := os.Getenv("appId")
	accessKeyId := os.Getenv("accessKeyId")
	accessKeySecret := os.Getenv("accessKeySecret")
	ctx := &dgctx.DgContext{TraceId: "123"}

	client := dgkdxf.NewClient(&dgkdxf.ClientConfig{
		AppId:           appId,
		Host:            host,
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
	})

	rt, err := client.GetAsrResult(ctx, os.Getenv("orderId"))
	if err != nil {
		panic(err)
	}
	dglogger.Infof(ctx, "rt: %s", utils.MustConvertBeanToJsonString(rt))
}

func TestConvert2Subtitles(t *testing.T) {
	orderResultBytes, _ := os.ReadFile("test.json")
	orderResult, _ := utils.ConvertJsonBytesToBean[dgkdxf.OrderResult](orderResultBytes)
	subtitlesList := orderResult.Convert2Subtitles()
	ctx := &dgctx.DgContext{TraceId: "123"}
	dglogger.Infof(ctx, "rt: %s", utils.MustConvertBeanToJsonString(subtitlesList))
	err := dgkdxf.ConvertSubtitles2SrtFormat(subtitlesList, "test.srt")
	if err != nil {
		panic(err)
	}
}
