package Login

import (
	wxClient "wechatwebapi/Cilent"
	"wechatwebapi/comm"
)

func Get62Data(Wxid string) string {
	D, err := comm.GetLoginata(Wxid)
	if err != nil {
		return err.Error()
	}
	return wxClient.Get62Data(D.Deviceid_str)
}
