package Group

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"strings"
	wxClient "wechatwebapi/Cilent"
	"wechatwebapi/Cilent/mm"
	"wechatwebapi/comm"
)

func OperateChatRoomAdmin(Data OperateChatRoomAdminParam) wxClient.ResponseResult {
	D, err := comm.GetLoginData(Data.Wxid)
	if err != nil {
		return wxClient.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：%v", err.Error()),
			Data:    nil,
		}
	}

	Cgiurl := "/cgi-bin/micromsg-bin/addchatroomadmin"
	Cgi := 889

	if Data.Val == 2 {
		Cgiurl = "/cgi-bin/micromsg-bin/delchatroomadmin"
		Cgi = 259
	}

	if Data.Val == 3 {
		Cgiurl = "/cgi-bin/micromsg-bin/transferchatroomowner"
		Cgi = 990
	}

	TowxdsSplit := strings.Split(Data.ToWxids, ",")

	var Towxds []string

	if len(TowxdsSplit) >= 1 {
		for _, v := range TowxdsSplit {
			Towxds = append(Towxds, v)
		}
	}

	req := &mm.ChatRoomAdminRequest{
		BaseRequest: &mm.BaseRequest{
			SessionKey:    D.SessionKey,
			Uin:           proto.Uint32(D.Uin),
			DeviceId:      D.Deviceid_byte,
			ClientVersion: proto.Int32(int32(wxClient.WxClientVersion)),
			DeviceType:    wxClient.DeviceTypeByte,
			Scene:         proto.Uint32(0),
		},
		ChatRoomName: proto.String(Data.QID),
		UserNameList: Towxds,
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
		Cgiurl:        Cgiurl,
		Proxy:         D.Proxy,
		Encryption:    5,
		TwelveEncData: wxClient.PackSpecialCgiData{},
		PackData: wxClient.PackData{
			Reqdata:          reqData,
			Cgi:              Cgi,
			Uin:              D.Uin,
			Cookie:           D.Cookie,
			SessionKey:       D.SessionKey,
			EncryptType:      5,
			Loginecdhkey:     D.Loginecdhkey,
			Clientsessionkey: D.Clientsessionkey,
			UseCompress:      true,
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
	Response := mm.ChatRoomAdminResponse{}
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
		Success: true,
		Message: "成功",
		Data:    &Response,
	}

}
