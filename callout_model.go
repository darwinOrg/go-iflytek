package iflytek

type CnoReq struct {
	Cno string `json:"cno"`
}

type CnoDetailResp struct {
	RequestId string `json:"requestId"`
	Client    struct {
		Cno     string `json:"cno"`
		Name    string `json:"name"`
		BindTel string `json:"bindTel"`
		Active  int    `json:"active"` // 是否激活，0: 否；1: 是
		Status  int    `json:"status"` // 座席状态，0: 离线；1: 在线
	} `json:"client"`
}

type CalloutReq struct {
	Cno             string `json:"cno"`
	CustomerNumber  string `json:"customerNumber"`
	RequestUniqueId string `json:"requestUniqueId"`
}

type CalloutResp struct {
	Result struct {
		Cno             string `json:"cno"`
		CustomerNumber  string `json:"customerNumber"`
		RequestUniqueId string `json:"requestUniqueId"`
	} `json:"result"`
	RequestId string `json:"requestId"`
}

type OnlineReq struct {
	Cno      string `json:"cno"`
	BindType int32  `json:"bindType"` // 电话类型，1:电话；2:IP话机
	BindTel  string `json:"bindTel"`  // 绑定电话
}

type RequestIdResp struct {
	RequestId string `json:"requestId"`
}

type OfflineReq struct {
	Cno       string `json:"cno"`
	UnbindTel int32  `json:"unbindTel"` // 是否下线同时解绑电话，0:不解绑；1:解绑
}

type BindClientTelReq struct {
	Cno    string `json:"cno"`
	Tel    string `json:"tel"`
	IsBind int    `json:"isBind"` // 是否绑定 1: 是，0: 否
}

type UnbindClientTelReq struct {
	Cno string `json:"cno"`
}

type OfflineResp struct {
	RequestId string `json:"requestId"`
}

type ListCdrObsReq struct {
	HiddenType     int32  `json:"hiddenType"`     // 是否隐藏号码。 0: 不隐藏，1: 中间四位，2: 最后八位 3: 全部号码，4: 最后四位。
	Cno            string `json:"cno"`            // 座席号
	CustomerNumber string `json:"customerNumber"` // 客户号码
	Status         int32  `json:"status"`         // 接听状态 0: 全部 1: 客户未接听 2: 座席未接听 3: 双方接听
}

type DownloadRecordFileReq struct {
	MainUniqueId string `binding:"required" json:"mainUniqueId"` // 通话记录唯一标识
	RecordSide   int32  `json:"recordSide"`                      // 不传递获取mp3格式录音，传递时获取wav格式录音。1：双轨录音客户侧，2：双轨录音座席侧，3：两侧合成录音
	RecordType   string `json:"recordType"`                      // "record": 通话录音，"voicemail": 留言。默认值为 "record"
}

type DownloadRecordFileResp struct {
	Filepath string `json:"filepath"` // 文件地址
	Filename string `json:"filename"` // 文件名称
}

type DownloadDetailRecordFileReq struct {
	MainUniqueId string `binding:"required" json:"mainUniqueId"` // 通话记录唯一标识
	UniqueId     string `binding:"required" json:"uniqueId"`     // 通话记录唯一标识
	RecordSide   int32  // 不传递获取mp3格式录音，传递时获取wav格式录音。1：双轨录音客户侧，2：双轨录音座席侧，3：两侧合成录音
}
