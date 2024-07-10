package Friend

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	wxClient "wechatwebapi/Cilent"
	"wechatwebapi/Cilent/mm"
	"wechatwebapi/bts"
	"wechatwebapi/comm"
)

func SetContactBlacklist(Data BlacklistParam) wxClient.ResponseResult {
	D, err := comm.GetLoginData(Data.Wxid)
	if err != nil {
		return wxClient.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：%v", err.Error()),
			Data:    nil,
		}
	}

	//先读取信息
	getContact := GetContact(GetContactDetailparameter{
		Wxid:     Data.Wxid,
		Towxids:  Data.ToWxids,
		ChatRoom: "",
	})

	if getContact.Code != 0 {
		return getContact
	}

	Contact := bts.GetContactResponse(getContact.Data)

	if len(Contact.ContactList) > 0 {
		modContact := Contact.ContactList[0]
		bit := uint32(0)
		if Data.Enable == 1 {
			bit = *(modContact.BitVal) | uint32(1<<3)
		} else {
			bit = *(modContact.BitVal) &^ uint32(1<<3)
		}
		ContactList := &mm.ModContact{
			UserName:        modContact.UserName,
			NickName:        modContact.NickName,
			PyInitial:       modContact.PyInitial,
			QuanPin:         modContact.QuanPin,
			Sex:             modContact.Sex,
			ImgBuf:          modContact.ImgBuf,
			BitMask:         modContact.BitMask,
			BitVal:          proto.Uint32(bit),
			ImgFlag:         modContact.ImgFlag,
			Remark:          modContact.Remark,
			RemarkPyinitial: modContact.RemarkPyinitial,
			RemarkQuanPin:   modContact.RemarkQuanPin,
			ContactType:     modContact.ContactType,
			ChatRoomNotify:  proto.Uint32(1),
			AddContactScene: modContact.AddContactScene,
			Extflag:         modContact.Extflag,
		}

		var cmdItems []*mm.CmdItem
		buffer, err := proto.Marshal(ContactList)
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

	return wxClient.ResponseResult{}
}
