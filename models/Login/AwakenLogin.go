package Login

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"strconv"
	"time"
	wxClient "wechatwebapi/Cilent"
	"wechatwebapi/Cilent/mm"
	"wechatwebapi/comm"
)

func AwakenLogin(Wxid string) wxClient.ResponseResult {
	D, err := comm.GetLoginata(Wxid)
	if err != nil {
		return wxClient.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：%v", err.Error()),
			Data:    nil,
		}
	}

	//初始化Mmtls
	_, MmtlsClient, err := comm.MmtlsInitialize(D.Proxy)
	if err != nil {
		return wxClient.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("MMTLS初始化失败：%v", err.Error()),
			Data:    nil,
		}
	}

	Autoauthkey := &mm.AutoAuthKey{}
	_ = proto.Unmarshal(D.Autoauthkey, Autoauthkey)

	req := &mm.PushLoginURLRequest{
		BaseRequest: &mm.BaseRequest{
			SessionKey:    D.Sessionkey,
			Uin:           proto.Uint32(D.Uin),
			DeviceId:      D.Deviceid_byte,
			ClientVersion: proto.Int32(int32(wxClient.WxClientVersion)),
			DeviceType:    wxClient.DeviceTypeByte,
			Scene:         proto.Uint32(0),
		},
		Autoauthticket: proto.String(""),
		Autoauthkey: &mm.SKBuiltinBufferT{
			ILen:   proto.Uint32(uint32(D.Autoauthkeylen)),
			Buffer: D.Autoauthkey,
		},
		ClientId:   proto.String("iPad-Push-" + strconv.Itoa(int(time.Now().Unix())) + ".110141"),
		Devicename: proto.String(D.DeviceName),
		Opcode:     proto.Int32(3),
		RandomEncryKey: &mm.SKBuiltinBufferT{
			ILen:   proto.Uint32(uint32(len(D.Sessionkey))),
			Buffer: D.Sessionkey,
		},
		Username: proto.String(D.Wxid),
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

	//开始发包
	protobufData, cookie, errType, err := comm.SendRequest(comm.SendPostData{
		Ip:         D.Mmtlsip,
		Host:       D.MmtlsHost,
		Cgiurl:     "/cgi-bin/micromsg-bin/pushloginurl",
		Proxy:      D.Proxy,
		Encryption: 12,
		TwelveEncData: wxClient.PackSpecialCgiData{
			Reqdata:                    reqData,
			Cgi:                        654,
			Encrypttype:                12,
			Extenddata:                 []byte{},
			Uin:                        D.Uin,
			Cookies:                    D.Cooike,
			ClientVersion:              wxClient.WxClientVersion,
			HybridEcdhPrivkey:          D.HybridEcdhPrivkey,
			HybridEcdhPubkey:           D.HybridEcdhPubkey,
			HybridEcdhInitServerPubKey: D.HybridEcdhInitServerPubKey,
		},
	}, MmtlsClient)

	if err != nil {
		return wxClient.ResponseResult{
			Code:    errType,
			Success: false,
			Message: err.Error(),
			Data:    nil,
		}
	}

	//解包
	PushLoginURLResponse := mm.PushLoginURLResponse{}
	err = proto.Unmarshal(protobufData, &PushLoginURLResponse)
	if err != nil {
		return wxClient.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("反序列化失败：%v", err.Error()),
			Data:    nil,
		}
	}

	//保存redis
	err = comm.CreateLoginData(comm.LoginData{
		Uuid:                       PushLoginURLResponse.GetUuid(),
		Aeskey:                     D.Sessionkey,
		NotifyKey:                  PushLoginURLResponse.GetNotifyKey().GetBuffer(),
		Deviceid_str:               D.Deviceid_str,
		Deviceid_byte:              D.Deviceid_byte,
		DeviceName:                 D.DeviceName,
		HybridEcdhPrivkey:          D.HybridEcdhPrivkey,
		HybridEcdhPubkey:           D.HybridEcdhPubkey,
		HybridEcdhInitServerPubKey: D.HybridEcdhInitServerPubKey,
		Cooike:                     cookie,
		MmtlsKey:                   MmtlsClient,
	}, "", 300)

	if err != nil {
		return wxClient.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("Redis ERROR：%v", err.Error()),
			Data:    nil,
		}
	}

	return wxClient.ResponseResult{
		Code:    0,
		Success: true,
		Message: "成功",
		Data:    PushLoginURLResponse,
	}
}
