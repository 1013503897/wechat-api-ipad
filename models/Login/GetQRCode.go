package Login

import (
	"encoding/hex"
	"fmt"
	"github.com/golang/protobuf/proto"
	"time"
	wxClient "wechatwebapi/Cilent"
	"wechatwebapi/Cilent/mm"
	"wechatwebapi/comm"
	"wechatwebapi/models"
)

type GetQRReq struct {
	Proxy      models.ProxyInfo
	DeviceID   string
	DeviceName string
}

type GetQRRes struct {
	QrBase64    string
	Uuid        string
	QrUrl       string
	ExpiredTime string
}

func GetQRCODE(DeviceID, DeviceName string, Proxy models.ProxyInfo) wxClient.ResponseResult {
	//初始化Mmtls
	_, MmtlsClient, err := comm.MmtlsInitialize(Proxy)
	if err != nil {
		return wxClient.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("MMTLS初始化失败：%v", err.Error()),
			Data:    nil,
		}
	}

	aesKey := []byte(wxClient.RandSeq(16)) //获取随机密钥
	deviceId := wxClient.CreateDeviceId(DeviceID)
	deviceIdByte, _ := hex.DecodeString(deviceId)

	HybridEcdhInitServerPubKey, HybridEcdhPrivKey, HybridEcdhPubKey := wxClient.HybridEcdhInit()

	req := &mm.GetLoginQRCodeRequest{
		BaseRequest: &mm.BaseRequest{
			SessionKey:    aesKey,
			Uin:           proto.Uint32(0),
			DeviceId:      deviceIdByte,
			ClientVersion: proto.Int32(int32(wxClient.WxClientVersion)),
			DeviceType:    wxClient.DeviceTypeByte,
			Scene:         proto.Uint32(0),
		},
		RandomEncryKey: &mm.SKBuiltinBufferT{
			ILen:   proto.Uint32(uint32(len(aesKey))),
			Buffer: aesKey,
		},
		Opcode:           proto.Uint32(0),
		MsgContextPubKey: nil,
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

	//开始请求发包
	protobufData, cookie, errType, err := comm.SendRequest(comm.SendPostData{
		Ip:         wxClient.MmtlsIp,
		Host:       wxClient.MmtlsHost,
		Cgiurl:     "/cgi-bin/micromsg-bin/getloginqrcode",
		Proxy:      Proxy,
		Encryption: 12,
		TwelveEncData: wxClient.PackSpecialCgiData{
			Reqdata:                    reqData,
			Cgi:                        501,
			Encrypttype:                12,
			Extenddata:                 []byte{},
			Uin:                        0,
			Cookies:                    []byte{},
			ClientVersion:              wxClient.WxClientVersion,
			HybridEcdhPrivkey:          HybridEcdhPrivKey,
			HybridEcdhPubkey:           HybridEcdhPubKey,
			HybridEcdhInitServerPubKey: HybridEcdhInitServerPubKey,
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

	getLoginQRRes := mm.GetLoginQRCodeResponse{}

	err = proto.Unmarshal(protobufData, &getLoginQRRes)

	if err != nil {
		return wxClient.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("反序列化失败：%v", err.Error()),
			Data:    nil,
		}
	}

	if getLoginQRRes.GetBaseResponse().GetRet() == 0 {
		uuid := getLoginQRRes.GetUuid()
		//保存redis
		err := comm.CreateLoginData(comm.LoginData{
			Uuid:                       uuid,
			Aeskey:                     aesKey,
			NotifyKey:                  getLoginQRRes.GetNotifyKey().GetBuffer(),
			Deviceid_str:               deviceId,
			Deviceid_byte:              deviceIdByte,
			DeviceName:                 DeviceName,
			HybridEcdhPrivkey:          HybridEcdhPrivKey,
			HybridEcdhPubkey:           HybridEcdhPubKey,
			HybridEcdhInitServerPubKey: HybridEcdhInitServerPubKey,
			Cooike:                     cookie,
			Proxy:                      Proxy,
			MmtlsKey:                   MmtlsClient,
		}, "", 300)

		if err == nil {
			return wxClient.ResponseResult{
				Code:    1,
				Success: true,
				Message: "成功",
				Data: GetQRRes{
					"",
					uuid,
					"http://weixin.qq.com/x/" + uuid,
					time.Unix(int64(getLoginQRRes.GetExpiredTime()), 0).Format("2006-01-02 15:04:05"),
				},
			}
		}
	}

	return wxClient.ResponseResult{
		Code:    -0,
		Success: false,
		Message: "未知的错误",
		Data:    getLoginQRRes,
	}
}
