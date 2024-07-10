package Login

import (
	"fmt"
	wxClient "wechatwebapi/Cilent"
	"wechatwebapi/comm"
)

func CacheInfo(Wxid string) wxClient.ResponseResult {
	D, err := comm.GetLoginData(Wxid)
	if err != nil {
		return wxClient.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：%v", err.Error()),
			Data:    nil,
		}
	}

	return wxClient.ResponseResult{
		Code:    1,
		Success: true,
		Message: "成功",
		Data:    D,
	}
}
