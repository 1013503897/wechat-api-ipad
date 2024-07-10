package Favor

import (
	"encoding/base64"
	"fmt"
	"github.com/golang/protobuf/proto"
	wxClient "wechatwebapi/Cilent"
	"wechatwebapi/Cilent/mm"
	"wechatwebapi/comm"
)

type SyncParam struct {
	Wxid   string
	Keybuf string
}

type SyncResponse struct {
	Ret    int32
	List   []mm.AddFavItem
	KeyBuf mm.SKBuiltinBufferT
}

func Sync(Data SyncParam) wxClient.ResponseResult {
	D, err := comm.GetLoginData(Data.Wxid)
	if err != nil {
		return wxClient.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：%v", err.Error()),
			Data:    nil,
		}
	}

	var KeyBuf mm.SKBuiltinBufferT

	if Data.Keybuf != "" {
		key, _ := base64.StdEncoding.DecodeString(Data.Keybuf)
		KeyBuf.Buffer = key
		KeyBuf.ILen = proto.Uint32(uint32(len(key)))
	}

	req := &mm.FavSyncRequest{
		Selector: proto.Uint32(1),
		KeyBuf:   &KeyBuf,
	}

	reqData, err := proto.Marshal(req)

	if err != nil {
		return wxClient.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("系统异常：%v", err.Error()),
			Data:    nil,
		}
	}

	//发包
	protobufData, _, errType, err := comm.SendRequest(comm.SendPostData{
		Ip:            D.Mmtlsip,
		Host:          D.MmtlsHost,
		Cgiurl:        "/cgi-bin/micromsg-bin/favsync",
		Proxy:         D.Proxy,
		Encryption:    5,
		TwelveEncData: wxClient.PackSpecialCgiData{},
		PackData: wxClient.PackData{
			Reqdata:          reqData,
			Cgi:              400,
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
	Response := mm.FavSyncResponse{}
	err = proto.Unmarshal(protobufData, &Response)
	if err != nil {
		return wxClient.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("反序列化失败：%v", err.Error()),
			Data:    nil,
		}
	}

	var List []mm.AddFavItem

	for _, v := range Response.CmdList.List {
		if *v.CmdId == int32(mm.SyncCmdID_MM_FAV_SYNCCMD_ADDITEM) {
			var data mm.AddFavItem
			_ = proto.Unmarshal(v.CmdBuf.Buffer, &data)
			List = append(List, data)
		}
	}

	return wxClient.ResponseResult{
		Code:    0,
		Success: true,
		Message: "成功",
		Data: SyncResponse{
			Ret:    *Response.Ret,
			List:   List,
			KeyBuf: *Response.KeyBuf,
		},
	}

}
