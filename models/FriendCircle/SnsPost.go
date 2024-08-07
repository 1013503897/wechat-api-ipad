package FriendCircle

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"strings"
	"time"
	wxClient "wechatwebapi/Cilent"
	"wechatwebapi/Cilent/mm"
	"wechatwebapi/comm"
)

type Messagearameter struct {
	Wxid         string
	Content      string
	BlackList    string
	WithUserList string
}

func Messages(Data Messagearameter) wxClient.ResponseResult {
	D, err := comm.GetLoginData(Data.Wxid)
	if err != nil {
		return wxClient.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：%v", err.Error()),
			Data:    nil,
		}
	}

	B := strings.Split(Data.BlackList, ",")
	BS := make([]*mm.SKBuiltinStringT, len(B))

	if len(B) >= 1 {
		for i, v := range B {
			BS[i] = &mm.SKBuiltinStringT{
				String_: proto.String(v),
			}
		}
	}

	W := strings.Split(Data.WithUserList, ",")
	WS := make([]*mm.SKBuiltinStringT, len(W))

	if len(W) >= 1 {
		for i, v := range W {
			WS[i] = &mm.SKBuiltinStringT{
				String_: proto.String(v),
			}
		}
	}

	ccData := &mm.CryptoData{
		Version:     []byte("00000003"),
		Type:        proto.Uint32(1),
		EncryptData: wxClient.GetNewSpamData(D.DeviceidStr, D.DeviceName),
		Timestamp:   proto.Uint32(uint32(time.Now().Unix())),
		Unknown5:    proto.Uint32(5),
		Unknown6:    proto.Uint32(0),
	}
	ccDataseq, _ := proto.Marshal(ccData)
	WCExtInfo := &mm.WCExtInfo{
		CcData: &mm.SKBuiltinBufferT{
			ILen:   proto.Uint32(uint32(len(ccDataseq))),
			Buffer: ccDataseq,
		},
	}
	_, _ = proto.Marshal(WCExtInfo)

	req := &mm.SnsPostRequest{
		BaseRequest: &mm.BaseRequest{
			SessionKey:    D.SessionKey,
			Uin:           proto.Uint32(D.Uin),
			DeviceId:      D.Deviceid_byte,
			ClientVersion: proto.Int32(int32(wxClient.WxClientVersion)),
			DeviceType:    wxClient.DeviceTypeByte,
			Scene:         proto.Uint32(0),
		},
		ObjectDesc: &mm.SKBuiltinString_S{
			ILen:   proto.Uint32(uint32(len(Data.Content))),
			Buffer: proto.String(Data.Content),
		},
		WithUserListNum: proto.Uint32(uint32(len(W))),
		WithUserList:    WS,
		ClientId:        proto.String(fmt.Sprintf("sns_post_%v_%v_0", D.Wxid, time.Now().Unix())),
		BlackListNum:    proto.Uint32(uint32(len(B))),
		BlackList:       BS,
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
		Cgiurl:        "/cgi-bin/micromsg-bin/mmsnspost",
		Proxy:         D.Proxy,
		Encryption:    5,
		TwelveEncData: wxClient.PackSpecialCgiData{},
		PackData: wxClient.PackData{
			Reqdata:          reqData,
			Cgi:              209,
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
	SnsPostResponse := mm.SnsPostResponse{}
	err = proto.Unmarshal(protobufData, &SnsPostResponse)

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
		Data:    SnsPostResponse,
	}

}
