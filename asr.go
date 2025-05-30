package iflytek

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	dgcoll "github.com/darwinOrg/go-common/collection"
	dgctx "github.com/darwinOrg/go-common/context"
	dgerr "github.com/darwinOrg/go-common/enums/error"
	"github.com/darwinOrg/go-common/model"
	"github.com/darwinOrg/go-common/utils"
	dghttp "github.com/darwinOrg/go-httpclient"
	dglogger "github.com/darwinOrg/go-logger"
	"github.com/google/uuid"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf8"
)

const (
	orderFinishedStatus = 4
)

var subtitlesSeparators = []string{"，", "。", "？", "！", ",", ".", "?", "!"}

type AsrUploadResult struct {
	Code     string `json:"code"`
	DescInfo string `json:"descInfo"`
	Content  struct {
		OrderId          string `json:"orderId"`
		TaskEstimateTime int    `json:"taskEstimateTime"`
	} `json:"content"`
}

func (r *AsrUploadResult) String() string {
	j, err := json.Marshal(r)
	if err != nil {
		return err.Error()
	} else {
		return string(j)
	}
}

type AsrResult struct {
	Code     string `json:"code"`
	DescInfo string `json:"descInfo"`
	Content  struct {
		OrderInfo struct {
			OrderId          string `json:"orderId"`
			FailType         int    `json:"failType"`
			Status           int    `json:"status"`
			OriginalDuration int    `json:"originalDuration"`
			RealDuration     int    `json:"realDuration"`
			ExpireTime       int    `json:"expireTime"`
		}
		OrderResultString string       `json:"orderResult"`
		OrderResult       *OrderResult `json:"-"`
	} `json:"content"`
}

type OrderResult struct {
	Lattice []struct {
		Json1best string `json:"json_1best"`
	} `json:"lattice"`
}

type Json1best struct {
	St struct {
		Sc string `json:"sc"`
		Pa string `json:"pa"`
		Rt []struct {
			Ws []struct {
				Cw []struct {
					W  string `json:"w"`
					Wp string `json:"wp"`
					Wc string `json:"wc"`
				} `json:"cw"`
				Wb int `json:"wb"`
				We int `json:"we"`
			} `json:"ws"`
		} `json:"rt"`
		Bg string `json:"bg"`
		Ed string `json:"ed"`
		Rl string `json:"rl"`
	} `json:"st"`
}

func (r *AsrResult) String() string {
	j, err := json.Marshal(r)
	if err != nil {
		return err.Error()
	} else {
		return string(j)
	}
}

func (o *OrderResult) Convert2Subtitles() []*Subtitles {
	if len(o.Lattice) == 0 {
		return []*Subtitles{}
	}

	var subtitlesList []*Subtitles
	var subtitlesBuilder strings.Builder
	var subtitlesBegin int

	for _, lattice := range o.Lattice {
		if lattice.Json1best == "" {
			continue
		}
		json1best := utils.MustConvertJsonStringToBean[Json1best](lattice.Json1best)
		latticeBegin, _ := strconv.ParseInt(json1best.St.Bg, 10, 0)

		for _, rt := range json1best.St.Rt {
			for _, ws := range rt.Ws {
				if subtitlesBegin == 0 {
					subtitlesBegin = int(latticeBegin) + ws.Wb*10
				}

				for _, cw := range ws.Cw {
					word := strings.TrimSpace(cw.W)
					if utf8.RuneCountInString(word) == 1 && dgcoll.Contains(subtitlesSeparators, word) {
						subtitlesList = append(subtitlesList, &Subtitles{
							Begin:     subtitlesBegin,
							End:       int(latticeBegin) + ws.We*10,
							Separator: cw.W,
							Words:     subtitlesBuilder.String(),
						})
						subtitlesBuilder.Reset()
						subtitlesBegin = 0
					} else {
						subtitlesBuilder.WriteString(cw.W)
					}
				}
			}
		}
	}

	return subtitlesList
}

func (o *OrderResult) String() string {
	var totalContent strings.Builder
	for _, lattice := range o.Lattice {
		if lattice.Json1best == "" {
			continue
		}
		json1best := utils.MustConvertJsonStringToBean[Json1best](lattice.Json1best)
		var itemStr strings.Builder
		rl := json1best.St.Rl
		for _, rt := range json1best.St.Rt {
			for _, ws := range rt.Ws {
				for _, cw := range ws.Cw {
					itemStr.WriteString(cw.W)
				}
			}
		}
		totalContent.WriteString(fmt.Sprintf("发言人%s: %s\n", rl, itemStr.String()))
	}

	return totalContent.String()
}

type sizedReader struct {
	r        io.Reader
	readSize int64
}

func (sr *sizedReader) Read(p []byte) (n int, err error) {
	n, err = sr.r.Read(p)
	if err == nil {
		sr.readSize += int64(n)
	}
	return n, err
}

// AsrUpload 上传录音文件到科大讯飞 返回 orderId failedReason err
func (c *Client) AsrUpload(dc *dgctx.DgContext, filePath string, duration int64, fileSize int64, callbackUrl string) (*AsrUploadResult, error) {
	file, err := os.Open(filePath)
	if err != nil {
		dglogger.Errorf(dc, "sdk OpenFile err: %v", err)
		return nil, err
	}

	defer func() {
		deferErr := file.Close()
		if deferErr != nil {
			dglogger.Errorf(dc, "Upload defer file.Close err: %v", deferErr)
		}
	}()

	uploadFileName := filepath.Base(filePath)
	params := c.buildUploadParams(uploadFileName, fileSize, duration, callbackUrl)
	parameters := utils.FormUrlEncodedParams(params)
	signature := c.GenerateSignature(params)
	uploadUrl := c.Config.Host + "/v2/upload?" + parameters
	reader := &sizedReader{
		r: bufio.NewReaderSize(file, defaultBufferSize),
	}

	req, err := http.NewRequest(http.MethodPost, uploadUrl, reader)
	if err != nil {
		dglogger.Errorf(dc, "sdk Upload http.NewRequest err: %v", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header["signature"] = []string{signature}

	response, err := dghttp.Client11.DoRequestRaw(dc, req)
	if err != nil {
		dglogger.Errorf(dc, "sdk Upload Client.Do err: %v", err)
		return nil, err
	}

	dglogger.Infof(dc, "sdk Upload %s file success, uploaded bytes size: %d, file size is:%d,url %s", uploadFileName, reader.readSize, fileSize, uploadUrl)

	if response.StatusCode != http.StatusOK {
		dglogger.Errorf(dc, "sdk Upload http.Post statusCode: %d", response.StatusCode)
		return nil, errors.New("upload call asr-service failed")
	}

	ret, err := dghttp.ConvertResponse2Struct[AsrUploadResult](response)
	if err != nil {
		dglogger.Errorf(dc, "sdk Upload utils.ResToObj err: %v", err)
		return nil, err
	}

	if ret.Code != apiSuccessCode {
		dglogger.Errorf(dc, "sdk Upload asr-service failed: %s", ret.String())
		return ret, ApiNoSuccessErr
	}

	return ret, nil
}

// GetAsrResult 获取科大讯飞的识别结果 api结果内容,音频识别内容,失败原因,error
func (c *Client) GetAsrResult(ctx *dgctx.DgContext, orderId string) (*AsrResult, error) {
	params := c.buildGetResultParams(orderId)
	formUrlString := utils.FormUrlEncodedParams(params)
	signature := c.GenerateSignature(params)
	resultUrl := c.Config.Host + "/v2/getResult?" + formUrlString

	dghttp.SetHttpClient(ctx, dghttp.Client11)
	defer dghttp.SetHttpClient(ctx, nil)
	ret, err := dghttp.DoGetToStruct[AsrResult](ctx, resultUrl, nil, map[string]string{"signature": signature})
	if err != nil {
		dglogger.Errorf(ctx, "dghttp.DoGetToStruct error | resultUrl: %s | err: %v", resultUrl, err)
		return nil, err
	}
	if ret == nil {
		return nil, dgerr.SYSTEM_ERROR
	}

	if ret.Code != apiSuccessCode {
		dglogger.Errorf(ctx, "sdk GetResult asr-service failed: %s", ret.String())
		return ret, ApiNoSuccessErr
	}

	orderInfo := ret.Content.OrderInfo
	if orderInfo.FailType != 0 {
		dglogger.Errorf(ctx, "sdk GetResult asr-service failed: %s", ret.String())
		reason := "order failType: " + strconv.FormatInt(int64(orderInfo.FailType), 10)
		return ret, errors.New(reason)
	}

	dglogger.Infof(ctx, "sdk GetResult orderId: %s,orderStatus: %d", orderId, orderInfo.Status)
	// 订单已完成的时候,解析识别结果
	if orderInfo.Status == orderFinishedStatus {
		ret.Content.OrderResult, err = utils.ConvertJsonStringToBean[OrderResult](ret.Content.OrderResultString)
		if err != nil {
			dglogger.Errorf(ctx, "sdk GetResult json.Unmarshal orderResult err: %v", err)
			return ret, err
		}
	}

	return ret, nil
}

func (c *Client) buildUploadParams(filename string, filesize int64, duration int64, callbackUrl string) []*model.KeyValuePair[string, any] {
	params := []*model.KeyValuePair[string, any]{
		{
			Key:   "dateTime",
			Value: getDateTimeString(),
		},
		{
			Key:   "accessKeyId",
			Value: c.Config.AccessKeyId,
		},
		{
			Key:   "signatureRandom",
			Value: uuid.NewString(),
		},
		{
			Key:   "fileName",
			Value: filename,
		},
		{
			Key:   "fileSize",
			Value: filesize,
		},
		{
			Key:   "duration",
			Value: duration,
		},
		{
			Key:   "language",
			Value: "cn",
		},
		{
			Key:   "roleType",
			Value: 1,
		},
		{
			Key:   "roleNum",
			Value: 0,
		},
		{
			Key:   "languageType",
			Value: 1,
		},
		{
			Key:   "callbackUrl",
			Value: callbackUrl,
		},
	}

	sortParams(params)
	return params
}

func (c *Client) buildGetResultParams(orderId string) []*model.KeyValuePair[string, any] {
	params := []*model.KeyValuePair[string, any]{
		{
			Key:   "dateTime",
			Value: getDateTimeString(),
		},
		{
			Key:   "accessKeyId",
			Value: c.Config.AccessKeyId,
		},
		{
			Key:   "signatureRandom",
			Value: uuid.NewString(),
		},
		{
			Key:   "orderId",
			Value: orderId,
		},
	}

	sortParams(params)
	return params
}
