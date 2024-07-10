package Friend

import (
	"fmt"
	wxClient "wechatwebapi/Cilent"
	"wechatwebapi/Cilent/mm"
	"wechatwebapi/comm"

	"github.com/golang/protobuf/proto"
)

type GetContactStatusParameter struct {
	Wxid     string
	UserList []string
}

func GetContactStatus(Data GetContactStatusParameter) wxClient.ResponseResult {
	D, err := comm.GetLoginData(Data.Wxid)
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

	if len(Data.UserList) >= 1 {
		for _, v := range Data.UserList {
			UserNameList = append(UserNameList, &mm.SKBuiltinStringT{
				String_: proto.String(v),
			})
		}
	}

	req := &mm.GetContactRequest{
		BaseRequest: &mm.BaseRequest{
			SessionKey:    D.SessionKey,
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

	status := make(map[string]int)
	if *Response.BaseResponse.Ret == 0 {
		for i := 0; i < int(*Response.ContactCount); i++ {
			value := 1
			user := Response.ContactList[i]
			if Response.Ticket[i].Antispamticket != nil && *Response.Ticket[i].Antispamticket != "" {
				if user.BigHeadImgUrl != nil && *user.BigHeadImgUrl == "" {
					value = 0
				} else {
					value = -1
				}
			}
			status[*Response.ContactList[i].UserName.String_] = value
		}
	}

	return wxClient.ResponseResult{
		Code:    0,
		Success: true,
		Message: "成功",
		Data:    status,
	}

}
