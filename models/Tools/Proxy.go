package Tools

import (
	"fmt"
	wxClient "wechatwebapi/Cilent"
	"wechatwebapi/comm"
	"wechatwebapi/models"
)

type SetProxyParam struct {
	Wxid  string
	Proxy models.ProxyInfo
}

func SetProxy(Data SetProxyParam) wxClient.ResponseResult {
	D, err := comm.GetLoginData(Data.Wxid)
	if err != nil {
		return wxClient.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：%v", err.Error()),
			Data:    nil,
		}
	}

	D.Proxy = Data.Proxy

	//初始化Mmtls
	_, MmtlsClient, err := comm.MmtlsInitialize(D.Proxy)
	if err != nil {
		return wxClient.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("MMTLS初始化失败：%v", err.Error()),
			Data:    nil,
		}
	}

	D.MmtlsKey = MmtlsClient

	err = comm.CreateLoginData(*D, Data.Wxid, 0)

	if err != nil {
		return wxClient.ResponseResult{
			Code:    -8,
			Success: false,
			Message: "失败",
			Data:    err.Error(),
		}
	}

	return wxClient.ResponseResult{
		Code:    1,
		Success: true,
		Message: "成功",
		Data:    nil,
	}

}
