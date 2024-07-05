package Tools

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"strings"
	wxClient "wechatwebapi/Cilent"
	"wechatwebapi/Cilent/mm"
	"wechatwebapi/comm"
)

type GetContactParam struct {
	Wxid         string
	UserNameList string
}

func GetContact(Data GetContactParam) wxClient.ResponseResult {
	D, err := comm.GetLoginata(Data.Wxid)
	if err != nil {
		return wxClient.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：%v", err.Error()),
			Data:    nil,
		}
	}

	var UserNameList []*mm.SKBuiltinStringT
	var AntispamTicket []*mm.SKBuiltinStringT

	UserNameListSplit := strings.Split(Data.UserNameList, ",")

	if len(UserNameListSplit) >= 1 {
		for _, v := range UserNameListSplit {
			UserNameList = append(UserNameList, &mm.SKBuiltinStringT{
				String_: proto.String(v),
			})
		}
	}

	req := &mm.GetContactRequest{
		BaseRequest: &mm.BaseRequest{
			SessionKey:    D.Sessionkey,
			Uin:           proto.Uint32(D.Uin),
			DeviceId:      D.Deviceid_byte,
			ClientVersion: proto.Int32(int32(wxClient.WxClientVersion)),
			DeviceType:    wxClient.DeviceTypeByte,
			Scene:         proto.Uint32(0),
		},
		UserCount:           proto.Int32(int32(len(UserNameList))),
		UserNameList:        UserNameList,
		AntispamTicketCount: proto.Int32(int32(len(AntispamTicket))),
		AntispamTicket:      AntispamTicket,
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
		Cgiurl:        "/cgi-bin/micromsg-bin/getcontact",
		Proxy:         D.Proxy,
		Encryption:    5,
		TwelveEncData: wxClient.PackSpecialCgiData{},
		PackData: wxClient.PackData{
			Reqdata:          reqData,
			Cgi:              182,
			Uin:              D.Uin,
			Cookie:           D.Cooike,
			Sessionkey:       D.Sessionkey,
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
	Response := mm.GetContactResponse{}
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
