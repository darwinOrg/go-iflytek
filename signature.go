package dgkdxf

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/url"
)

func (c *Client) GenerateSignature(params map[string]any) string {
	baseString := formUrlEncodedParams(params)

	return GenerateSignatureWithStd(c.Config.AccessKeySecret, baseString)
}

// GenerateSignatureWithStd key+baseString 先sha1 然后 base64.StdEncoding
func GenerateSignatureWithStd(key string, baseString string) string {
	hash := hmac.New(sha1.New, []byte(key))
	hash.Write([]byte(baseString))
	bytes := hash.Sum(nil)

	return base64.StdEncoding.EncodeToString(bytes)
}

func formUrlEncodedParams(params map[string]any) string {
	paramStr := ""
	for k, v := range params {
		if k != "" {
			paramStr += url.QueryEscape(k) + "=" + url.QueryEscape(fmt.Sprintf("%v", v)) + "&"
		}
	}

	paramLen := len(paramStr)
	if paramLen > 0 {
		paramStr = paramStr[:paramLen-1]
	}

	return paramStr
}
