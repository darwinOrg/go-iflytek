package iflytek

import (
	"bufio"
	"errors"
	"fmt"
	dgctx "github.com/darwinOrg/go-common/context"
	"github.com/darwinOrg/go-common/model"
	"github.com/darwinOrg/go-common/utils"
	dghttp "github.com/darwinOrg/go-httpclient"
	dglogger "github.com/darwinOrg/go-logger"
	"io"
	"mime"
	"net/http"
	"os"
	"strings"
)

// DetailByCno 查看坐席详情
func (c *Client) DetailByCno(ctx *dgctx.DgContext, conReq *CnoReq) (*CnoDetailResp, error) {
	uri := c.buildDetailByCnoUri(conReq.Cno)
	dglogger.Infof(ctx, "DetailByCno buildDetailByCnoUri: %s", uri)

	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		dglogger.Errorf(ctx, "DetailByCno http.NewRequest err: %v", err)
		return nil, err
	}

	response, err := dghttp.Client2.DoRequestRaw(ctx, req)
	if err != nil {
		dglogger.Errorf(ctx, "DetailByCno dghttp.Client2.DoRequestRaw err: %v", err)
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		dglogger.Errorf(ctx, "DetailByCno dghttp.Client2.DoRequestRaw statusCode: %d", response.StatusCode)
		return nil, errors.New(fmt.Sprintf("call detail cno statusCode: %d", response.StatusCode))
	}

	bytes, err := read(response)
	if err != nil {
		dglogger.Errorf(ctx, "DetailByCno read response err: %v", err)
		return nil, err
	}

	dglogger.Infof(ctx, "res: %s", string(bytes))
	calloutResp, err := utils.ConvertJsonBytesToBean[CnoDetailResp](bytes)
	if err != nil {
		dglogger.Errorf(ctx, "DetailByCno ConvertJsonBytesToBean err: %v", err)
		return nil, err
	}

	return calloutResp, nil
}

// Callout 外呼
func (c *Client) Callout(ctx *dgctx.DgContext, calloutReq *CalloutReq) (*CalloutResp, error) {
	uri := c.buildPostUri("/cc/callout?")
	dglogger.Infof(ctx, "Callout buildPostUri: %s", uri)

	req, err := http.NewRequest(http.MethodPost, uri, strings.NewReader(utils.MustConvertBeanToJsonString(calloutReq)))
	if err != nil {
		dglogger.Errorf(ctx, "Callout http.NewRequest err: %v", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response, err := dghttp.Client2.DoRequestRaw(ctx, req)
	if err != nil {
		dglogger.Errorf(ctx, "Callout dghttp.Client2.DoRequestRaw err: %v", err)
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		dglogger.Errorf(ctx, "Callout dghttp.Client2.DoRequestRaw statusCode: %d", response.StatusCode)
		return nil, errors.New(fmt.Sprintf("call callout statusCode: %d", response.StatusCode))
	}

	bytes, err := read(response)
	if err != nil {
		dglogger.Errorf(ctx, "Callout read response err: %v", err)
		return nil, err
	}

	dglogger.Infof(ctx, "res: %s", string(bytes))
	calloutResp, err := utils.ConvertJsonBytesToBean[CalloutResp](bytes)
	if err != nil {
		dglogger.Errorf(ctx, "Callout ConvertJsonBytesToBean err: %v", err)
		return nil, err
	}

	return calloutResp, nil
}

// Cancel 外呼取消
func (c *Client) Cancel(ctx *dgctx.DgContext, conReq *CnoReq) (*RequestIdResp, error) {
	uri := c.buildPostUri("/cc/callout_cancel?")
	dglogger.Infof(ctx, "Cancel buildPostUri: %s", uri)

	req, err := http.NewRequest(http.MethodPost, uri, strings.NewReader(utils.MustConvertBeanToJsonString(conReq)))
	if err != nil {
		dglogger.Errorf(ctx, "Cancel http.NewRequest err: %v", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response, err := dghttp.Client2.DoRequestRaw(ctx, req)
	if err != nil {
		dglogger.Errorf(ctx, "Cancel dghttp.Client2.DoRequestRaw err: %v", err)
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		dglogger.Errorf(ctx, "Cancel dghttp.Client2.DoRequestRaw statusCode: %d", response.StatusCode)
		return nil, errors.New(fmt.Sprintf("call unlink statusCode: %d", response.StatusCode))
	}

	bytes, err := read(response)
	if err != nil {
		dglogger.Errorf(ctx, "Cancel read response err: %v", err)
		return nil, err
	}

	dglogger.Infof(ctx, "res: %s", string(bytes))
	cancelResp, err := utils.ConvertJsonBytesToBean[RequestIdResp](bytes)
	if err != nil {
		dglogger.Errorf(ctx, "Cancel ConvertJsonBytesToBean err: %v", err)
		return nil, err
	}

	return cancelResp, nil
}

// Unlink 挂机
func (c *Client) Unlink(ctx *dgctx.DgContext, conReq *CnoReq) (*RequestIdResp, error) {
	uri := c.buildPostUri("/cc/unlink?")
	dglogger.Infof(ctx, "Unlink buildPostUri: %s", uri)

	req, err := http.NewRequest(http.MethodPost, uri, strings.NewReader(utils.MustConvertBeanToJsonString(conReq)))
	if err != nil {
		dglogger.Errorf(ctx, "Unlink http.NewRequest err: %v", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response, err := dghttp.Client2.DoRequestRaw(ctx, req)
	if err != nil {
		dglogger.Errorf(ctx, "Unlink dghttp.Client2.DoRequestRaw err: %v", err)
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		dglogger.Errorf(ctx, "Unlink dghttp.Client2.DoRequestRaw statusCode: %d", response.StatusCode)
		return nil, errors.New(fmt.Sprintf("call unlink statusCode: %d", response.StatusCode))
	}

	bytes, err := read(response)
	if err != nil {
		dglogger.Errorf(ctx, "Unlink read response err: %v", err)
		return nil, err
	}

	dglogger.Infof(ctx, "res: %s", string(bytes))
	unlinkResp, err := utils.ConvertJsonBytesToBean[RequestIdResp](bytes)
	if err != nil {
		dglogger.Errorf(ctx, "Unlink ConvertJsonBytesToBean err: %v", err)
		return nil, err
	}

	return unlinkResp, nil
}

// Online 上线
func (c *Client) Online(ctx *dgctx.DgContext, onlineReq *OnlineReq) (*RequestIdResp, error) {
	uri := c.buildPostUri("/cc/online?")
	dglogger.Infof(ctx, "Online BuildOnlineUri: %s", uri)

	req, err := http.NewRequest(http.MethodPost, uri, strings.NewReader(utils.MustConvertBeanToJsonString(onlineReq)))
	if err != nil {
		dglogger.Errorf(ctx, "Online http.NewRequest err: %v", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response, err := dghttp.Client2.DoRequestRaw(ctx, req)
	if err != nil {
		dglogger.Errorf(ctx, "Online dghttp.Client2.DoRequestRaw err: %v", err)
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		dglogger.Errorf(ctx, "Online dghttp.Client2.DoRequestRaw statusCode: %d", response.StatusCode)
		return nil, errors.New(fmt.Sprintf("call online statusCode: %d", response.StatusCode))
	}

	bytes, err := read(response)
	if err != nil {
		dglogger.Errorf(ctx, "Online read response err: %v", err)
		return nil, errors.New(fmt.Sprintf("call online statusCode: %d", response.StatusCode))
	}

	dglogger.Infof(ctx, "res: %s", string(bytes))
	onlineResp, err := utils.ConvertJsonBytesToBean[RequestIdResp](bytes)
	if err != nil {
		dglogger.Errorf(ctx, "Online ConvertJsonBytesToBean err: %v", err)
		return nil, err
	}

	return onlineResp, nil
}

// Offline 下线
func (c *Client) Offline(ctx *dgctx.DgContext, offlineReq *OfflineReq) (*OfflineResp, error) {
	uri := c.buildPostUri("/cc/offline?")
	dglogger.Infof(ctx, "Offline buildPostUri: %s", uri)

	req, err := http.NewRequest(http.MethodPost, uri, strings.NewReader(utils.MustConvertBeanToJsonString(offlineReq)))
	if err != nil {
		dglogger.Errorf(ctx, "Offline http.NewRequest err: %v", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response, err := dghttp.Client2.DoRequestRaw(ctx, req)
	if err != nil {
		dglogger.Errorf(ctx, "Offline dghttp.Client2.DoRequestRaw err: %v", err)
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		dglogger.Errorf(ctx, "Offline dghttp.Client2.DoRequestRaw statusCode: %d", response.StatusCode)
		return nil, errors.New(fmt.Sprintf("call offline statusCode: %d", response.StatusCode))
	}

	bytes, err := read(response)
	if err != nil {
		dglogger.Errorf(ctx, "Offline read response err: %v", err)
		return nil, errors.New(fmt.Sprintf("call offline statusCode: %d", response.StatusCode))
	}

	dglogger.Infof(ctx, "res: %s", string(bytes))
	offlineResp, err := utils.ConvertJsonBytesToBean[OfflineResp](bytes)
	if err != nil {
		dglogger.Errorf(ctx, "Offline ConvertJsonBytesToBean err: %v", err)
		return nil, err
	}

	return offlineResp, nil
}

// ListCdrObs 查询外呼通话记录列表
func (c *Client) ListCdrObs(ctx *dgctx.DgContext, listCdrObsReq *ListCdrObsReq) {
	uri := c.buildListCdrObsUri(listCdrObsReq)
	dglogger.Infof(ctx, "ListCdrObs buildListCdrObsUri: %s", uri)

	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		dglogger.Errorf(ctx, "ListCdrObs http.NewRequest err: %v", err)
		return
	}

	response, err := dghttp.Client2.DoRequestRaw(ctx, req)
	if err != nil {
		dglogger.Errorf(ctx, "ListCdrObs dghttp.Client2.DoRequestRaw err: %v", err)
		return
	}

	if response.StatusCode != http.StatusOK {
		dglogger.Errorf(ctx, "ListCdrObs dghttp.Client2.DoRequestRaw statusCode: %d", response.StatusCode)
		return
	}

	bytes, err := read(response)
	if err != nil {
		dglogger.Errorf(ctx, "ListCdrObs read response err: %v", err)
		return
	}

	dglogger.Infof(ctx, "res: %s", string(bytes))
}

// DownloadRecordFile 下载通话详情录音文件
func (c *Client) DownloadRecordFile(ctx *dgctx.DgContext, downloadRecordFileReq *DownloadRecordFileReq) (*DownloadRecordFileResp, error) {
	uri := c.buildDownloadRecordFileUri(downloadRecordFileReq)
	dglogger.Infof(ctx, "DownloadRecordFile buildDownloadRecordFileUri: %s", uri)

	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		dglogger.Errorf(ctx, "DownloadRecordFile http.NewRequest err: %v", err)
		return nil, err
	}

	response, err := dghttp.Client11.DoRequestRaw(ctx, req)
	if err != nil {
		dglogger.Errorf(ctx, "DownloadRecordFile dghttp.Client11.DoRequestRaw err: %v", err)
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		dglogger.Errorf(ctx, "DownloadRecordFile dghttp.Client11.DoRequestRaw statusCode: %d", response.StatusCode)
		return nil, errors.New(fmt.Sprintf("call download_record_file statusCode: %d", response.StatusCode))
	}

	disposition := response.Header.Get("Content-Disposition")
	fileName := ""
	if disposition != "" {
		_, params, err := mime.ParseMediaType(disposition)
		if err != nil {
			dglogger.Errorf(ctx, "DownloadRecordFile ParseMediaType err: %v", err)

			return nil, err
		}
		var ok bool
		fileName, ok = params["filename"]
		if !ok {
			dglogger.Errorf(ctx, "DownloadRecordFile filename not exist")
			return nil, errors.New("DownloadRecordFile filename not exist")
		}
	}

	defer response.Body.Close()

	filepath := "/tmp/" + fileName
	file, err := os.Create(filepath)
	if err != nil {
		dglogger.Errorf(ctx, "DownloadRecordFile os.Create err: %v", err)
		return nil, err
	}
	defer file.Close()

	writer := bufio.NewWriterSize(file, defaultBufferSize)
	defer writer.Flush()

	_, err = io.Copy(writer, bufio.NewReaderSize(response.Body, defaultBufferSize))
	if err != nil {
		dglogger.Errorf(ctx, "DownloadRecordFile io.Copy err: %v", err)
		return nil, err
	}

	return &DownloadRecordFileResp{Filepath: filepath, Filename: fileName}, nil
}

// BindClientTel 绑定座席电话
func (c *Client) BindClientTel(ctx *dgctx.DgContext, bindReq *BindClientTelReq) error {
	uri := c.buildPostUri("/cc/bind_client_tel?")
	dglogger.Infof(ctx, "BindClientTel buildPostUri: %s", uri)
	dghttp.SetHttpClient(ctx, dghttp.Client2)
	defer dghttp.SetHttpClient(ctx, nil)
	resp, err := dghttp.DoPostJsonToStruct[KdxfResponse](ctx, uri, bindReq, nil)
	if err != nil {
		dglogger.Errorf(ctx, "BindClientTel[%+v] do post err: %v", bindReq, err)
		return err
	}
	if resp.Error.Message != "" {
		return errors.New(resp.Error.Message)
	}
	return nil
}

// UnbindClientTel 解绑座席电话
func (c *Client) UnbindClientTel(ctx *dgctx.DgContext, unbindReq *UnbindClientTelReq) error {
	uri := c.buildPostUri("/cc/unbind_client_tel?")
	dglogger.Infof(ctx, "UnbindClientTel buildPostUri: %s", uri)
	dghttp.SetHttpClient(ctx, dghttp.Client2)
	defer dghttp.SetHttpClient(ctx, nil)
	resp, err := dghttp.DoPostJsonToStruct[KdxfResponse](ctx, uri, unbindReq, nil)
	if err != nil {
		dglogger.Errorf(ctx, "UnbindClientTel[%+v] do post err: %v", unbindReq, err)
		return err
	}
	if resp.Error.Message != "" {
		return errors.New(resp.Error.Message)
	}
	return nil
}

func (c *Client) buildPostUri(callUrl string) string {
	params := c.buildCommonParams()
	after := cutPrefix(c.Config.Host, "https://")
	urlPrefix := fmt.Sprintf("%s%s%s", http.MethodPost, after, callUrl)
	signature := c.GenerateSignatureWithUrlPrefix(urlPrefix, params)

	params = append(params, &model.KeyValuePair[string, any]{Key: "Signature", Value: signature})
	parameters := utils.FormUrlEncodedParams(params)

	return c.Config.Host + callUrl + parameters
}

func (c *Client) buildDetailByCnoUri(cno string) string {
	params := c.buildCommonParams()
	params = append(params, &model.KeyValuePair[string, any]{Key: "cno", Value: cno})
	sortParams(params)

	callUrl := "/cc/describe_client?"
	after := cutPrefix(c.Config.Host, "https://")
	urlPrefix := fmt.Sprintf("%s%s%s", http.MethodGet, after, callUrl)
	signature := c.GenerateSignatureWithUrlPrefix(urlPrefix, params)

	params = append(params, &model.KeyValuePair[string, any]{Key: "Signature", Value: signature})
	parameters := utils.FormUrlEncodedParams(params)

	return c.Config.Host + callUrl + parameters
}

func (c *Client) buildListCdrObsUri(listCdrObsReq *ListCdrObsReq) string {
	params := c.buildCommonParams()
	params = append(params, &model.KeyValuePair[string, any]{Key: "hiddenType", Value: listCdrObsReq.HiddenType})
	params = append(params, &model.KeyValuePair[string, any]{Key: "customerNumber", Value: listCdrObsReq.CustomerNumber})
	params = append(params, &model.KeyValuePair[string, any]{Key: "cno", Value: listCdrObsReq.Cno})
	params = append(params, &model.KeyValuePair[string, any]{Key: "status", Value: listCdrObsReq.Status})
	sortParams(params)

	callUrl := "/cc/list_cdr_obs?"
	after := cutPrefix(c.Config.Host, "https://")
	urlPrefix := fmt.Sprintf("%s%s%s", http.MethodGet, after, callUrl)
	signature := c.GenerateSignatureWithUrlPrefix(urlPrefix, params)

	params = append(params, &model.KeyValuePair[string, any]{Key: "Signature", Value: signature})
	parameters := utils.FormUrlEncodedParams(params)

	return c.Config.Host + callUrl + parameters
}

func (c *Client) buildDownloadRecordFileUri(downloadRecordFileReq *DownloadRecordFileReq) string {
	params := c.buildCommonParams()
	params = append(params, &model.KeyValuePair[string, any]{Key: "mainUniqueId", Value: downloadRecordFileReq.MainUniqueId})
	if downloadRecordFileReq.RecordSide > 0 {
		params = append(params, &model.KeyValuePair[string, any]{Key: "recordSide", Value: downloadRecordFileReq.RecordSide})
	}
	params = append(params, &model.KeyValuePair[string, any]{Key: "recordType", Value: downloadRecordFileReq.RecordType})
	sortParams(params)

	callUrl := "/cc/download_record_file?"
	after := cutPrefix(c.Config.Host, "https://")
	urlPrefix := fmt.Sprintf("%s%s%s", http.MethodGet, after, callUrl)
	signature := c.GenerateSignatureWithUrlPrefix(urlPrefix, params)

	params = append(params, &model.KeyValuePair[string, any]{Key: "Signature", Value: signature})
	parameters := utils.FormUrlEncodedParams(params)

	return c.Config.Host + callUrl + parameters
}

func read(resp *http.Response) ([]byte, error) {
	if resp == nil {
		return nil, nil
	}
	if resp.Body == nil {
		return nil, nil
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func cutPrefix(s, prefix string) string {
	if !strings.HasPrefix(s, prefix) {
		return s
	}
	return s[len(prefix):]
}
