package dgkdxf_test

import (
	dgctx "github.com/darwinOrg/go-common/context"
	dgkdxf "github.com/darwinOrg/go-iflytek"
	dglogger "github.com/darwinOrg/go-logger"
	"os"
	"testing"
)

func TestBuildAstUri(t *testing.T) {
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
	uri := client.BuildAstUri(ctx, &dgkdxf.AstParamConfig{
		Lang:           "cn",
		Codec:          "pcm_s16le",
		AudioEncode:    "pcm",
		Samplerate:     "16000",
		RoleType:       dgkdxf.RoleTypeOpen,
		ContextId:      "",
		FeatureIds:     nil,
		HotWordId:      "",
		SourceInfo:     "",
		FilePath:       "",
		ResultFilePath: "",
	})
	dglogger.Infof(ctx, "uri: %s", uri)
}
