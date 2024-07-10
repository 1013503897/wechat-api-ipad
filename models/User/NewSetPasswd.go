package User

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/golang/protobuf/proto"
	wxClient "wechatwebapi/Cilent"
	"wechatwebapi/Cilent/mm"
	"wechatwebapi/comm"
)

type NewSetPasswdParam struct {
	Wxid     string
	Password string
	Ticket   string
}

func NewSetPasswd(Data NewSetPasswdParam) wxClient.ResponseResult {
	D, err := comm.GetLoginData(Data.Wxid)
	if err != nil {
		return wxClient.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：%v", err.Error()),
			Data:    nil,
		}
	}

	req := &mm.SetPwdRequest{
		BaseRequest: &mm.BaseRequest{
			SessionKey:    D.SessionKey,
			Uin:           proto.Uint32(D.Uin),
			DeviceId:      D.Deviceid_byte,
			ClientVersion: proto.Int32(int32(D.ClientVersion)),
			DeviceType:    []byte(D.DeviceType),
			Scene:         proto.Uint32(0),
		},
		Password: proto.String(comm.MD5ToLower(Data.Password)),
		Ticket:   proto.String(Data.Ticket),
		AutoAuthKey: &mm.SKBuiltinBufferT{
			ILen:   proto.Uint32(uint32(len(D.Autoauthkey))),
			Buffer: D.Autoauthkey,
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

	S2801, _ := hex.DecodeString("2801")

	reqDataA := new(bytes.Buffer)
	reqDataA.Write(reqData)

	if len(D.Deviceid_byte) <= 16 {
		reqDataA.Write(S2801)
	}

	//发包
	protobufData, _, errType, err := comm.SendRequest(comm.SendPostData{
		Ip:            D.Mmtlsip,
		Cgiurl:        "/cgi-bin/micromsg-bin/newsetpasswd",
		Proxy:         D.Proxy,
		Encryption:    5,
		TwelveEncData: wxClient.PackSpecialCgiData{},
		PackData: wxClient.PackData{
			Reqdata:          reqDataA.Bytes(),
			Cgi:              383,
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
	Response := mm.SetPwdResponse{}
	err = proto.Unmarshal(protobufData, &Response)
	if err != nil {
		return wxClient.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("反序列化失败：%v", err.Error()),
			Data:    nil,
		}
	}

	//更新成功就保存autoAuthKey
	if Response.AutoAuthKey != nil && len(Response.AutoAuthKey.Buffer) > 10 {
		D.Autoauthkey = Response.AutoAuthKey.Buffer
		err = comm.CreateLoginData(*D, D.Wxid, 0)
		if err != nil {
			return wxClient.ResponseResult{
				Code:    -8,
				Success: false,
				Message: fmt.Sprintf("AutoAuthKey保存失败：%v", err.Error()),
				Data:    nil,
			}
		}
	}

	return wxClient.ResponseResult{
		Code:    0,
		Success: true,
		Message: "成功",
		Data:    &Response,
	}

}
