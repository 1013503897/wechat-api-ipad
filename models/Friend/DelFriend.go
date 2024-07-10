package Friend

import (
	"fmt"
	wxClient "wechatwebapi/Cilent"
	"wechatwebapi/Cilent/mm"
	"wechatwebapi/comm"

	"github.com/golang/protobuf/proto"
)

type DelFriendParam struct {
	Wxid     string
	UserName string
}

func DelFriend(Data DelFriendParam) wxClient.ResponseResult {
	D, err := comm.GetLoginData(Data.Wxid)
	if err != nil {
		return wxClient.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：%v", err.Error()),
			Data:    nil,
		}
	}
	delContact := &mm.DelContact{
		DeleteContactScene: proto.Uint32(uint32(0)),
		UserName: &mm.SKBuiltinStringT{
			String_: proto.String(Data.UserName),
		},
	}

	//序列化
	delContactSerialized, _ := proto.Marshal(delContact)

	var CmdItem []*mm.CmdItem

	CmdItem = append(CmdItem, &mm.CmdItem{
		CmdId: proto.Int32(4),
		CmdBuf: &mm.SKBuiltinBufferT{
			ILen:   proto.Uint32(uint32(len(delContactSerialized))),
			Buffer: delContactSerialized,
		},
	})

	req := &mm.OpLogRequest{
		Cmd: &mm.CmdList{
			Count: proto.Uint32(uint32(len(CmdItem))),
			List:  CmdItem,
		},
	}

	//序列化
	reqData, _ := proto.Marshal(req)

	//发包
	protobufData, _, errType, err := comm.SendRequest(comm.SendPostData{
		Ip:            D.Mmtlsip,
		Cgiurl:        "/cgi-bin/micromsg-bin/oplog",
		Proxy:         D.Proxy,
		Encryption:    5,
		TwelveEncData: wxClient.PackSpecialCgiData{},
		PackData: wxClient.PackData{
			Reqdata:          reqData,
			Cgi:              681,
			Uin:              D.Uin,
			Cookie:           D.Cookie,
			SessionKey:       D.SessionKey,
			EncryptType:      5,
			Loginecdhkey:     D.Loginecdhkey,
			Clientsessionkey: D.Clientsessionkey,
			UseCompress:      false,
		},
	}, D.MmtlsKey)

	if err != nil {
		return wxClient.ResponseResult{
			Code:    errType,
			Success: false,
			Message: err.Error(),
			Data:    nil,
		}
	}

	//解包
	Response := mm.OplogResponse{}
	err = proto.Unmarshal(protobufData, &Response)
	if err != nil {
		return wxClient.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("反序列化失败：%v", err.Error()),
			Data:    nil,
		}
	}

	return wxClient.ResponseResult{
		Code:    0,
		Success: false,
		Message: "成功",
		Data:    &Response,
	}
}
