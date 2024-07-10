package Login

import (
	wxClient "wechatwebapi/Cilent"
	"wechatwebapi/comm"
)

func Get62Data(Wxid string) string {
	D, err := comm.GetLoginData(Wxid)
	if err != nil {
		return err.Error()
	}
	return wxClient.Get62Data(D.DeviceidStr)
}
