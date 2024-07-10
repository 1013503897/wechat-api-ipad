package Group

import (
	"fmt"
	wxClient "wechatwebapi/Cilent"
	"wechatwebapi/Cilent/mm"
	"wechatwebapi/bts"
	"wechatwebapi/comm"
	"wechatwebapi/models/Tools"

	"github.com/golang/protobuf/proto"
)

func MoveToContract(Data MoveContractListParam) wxClient.ResponseResult {
	D, err := comm.GetLoginData(Data.Wxid)
	if err != nil {
		return wxClient.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：%v", err.Error()),
			Data:    nil,
		}
	}

	GetContact := Tools.GetContact(Tools.GetContactParam{
		Wxid:         Data.Wxid,
		UserNameList: Data.QID,
	})

	if GetContact.Data == nil {
		return wxClient.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("系统异常：%v", GetContact.Message),
			Data:    nil,
		}
	}

	Contact := bts.GetContactResponse(GetContact.Data)
	modContact := Contact.ContactList[0]
	bit := uint32(0)
	if Data.Val == 1 {
		bit = *(modContact.BitVal) | uint32(1<<0)
	} else {
		bit = *(modContact.BitVal) &^ uint32(1<<0)
	}

	ModContact := &mm.ModContact{
		UserName: &mm.SKBuiltinStringT{
			String_: proto.String(Data.QID),
		},
		NickName:  &mm.SKBuiltinStringT{},
		PyInitial: &mm.SKBuiltinStringT{},
		QuanPin:   &mm.SKBuiltinStringT{},
		Sex:       proto.Int32(0),
		ImgBuf:    &mm.SKBuiltinBufferT{},
		BitMask:   Contact.ContactList[0].BitMask,
		BitVal:    proto.Uint32(bit),
		ImgFlag:   proto.Uint32(0),
		Remark: &mm.SKBuiltinStringT{
			String_: Contact.ContactList[0].Remark.String_,
		},
		RemarkPyinitial: &mm.SKBuiltinStringT{
			String_: Contact.ContactList[0].RemarkPyinitial.String_,
		},
		RemarkQuanPin: &mm.SKBuiltinStringT{
			String_: Contact.ContactList[0].RemarkQuanPin.String_,
		},
		ContactType:     proto.Uint32(0),
		ChatRoomNotify:  proto.Uint32(1),
		AddContactScene: proto.Uint32(0),
		Extflag:         proto.Int32(0),
	}

	buffer, err := proto.Marshal(ModContact)
	if err != nil {
		return wxClient.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("系统异常：%v", err.Error()),
			Data:    nil,
		}
	}

	cmdItem := mm.CmdItem{
		CmdId: proto.Int32(2),
		CmdBuf: &mm.SKBuiltinBufferT{
			ILen:   proto.Uint32(uint32(len(buffer))),
			Buffer: buffer,
		},
	}

	var cmdItems []*mm.CmdItem
	cmdItems = append(cmdItems, &cmdItem)

	req := &mm.OpLogRequest{
		Cmd: &mm.CmdList{
			Count: proto.Uint32(uint32(len(cmdItems))),
			List:  cmdItems,
		},
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
	GetContactResponse := mm.OplogResponse{}
	err = proto.Unmarshal(protobufData, &GetContactResponse)

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
		Data:    GetContactResponse,
	}

}
