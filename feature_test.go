package iflytek_test

import (
	"encoding/base64"
	dgctx "github.com/darwinOrg/go-common/context"
	dgkdxf "github.com/darwinOrg/go-iflytek"
	"github.com/google/uuid"
	"os"
	"testing"
)

func TestRegisterFeature(t *testing.T) {
	audioData, _ := os.ReadFile("test.wav")
	audioDataStr := base64.StdEncoding.EncodeToString(audioData)
	req := &dgkdxf.RegisterFeatureRequest{
		AudioData: audioDataStr,
		AudioType: dgkdxf.AudioTypeRaw,
		Uid:       uuid.NewString(),
	}

	host := "https://office-api-personal-dx.iflyaisol.com"
	appId := os.Getenv("appId")
	accessKeyId := os.Getenv("accessKeyId")
	accessKeySecret := os.Getenv("accessKeySecret")
	ctx := &dgctx.DgContext{TraceId: uuid.NewString()}
	client := dgkdxf.NewClient(&dgkdxf.ClientConfig{
		AppId:           appId,
		Host:            host,
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
	})

	featureId, err := client.RegisterFeature(ctx, req)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(featureId)
	}
}
