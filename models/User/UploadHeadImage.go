package User

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/golang/protobuf/proto"
	"strings"
	wxClient "wechatwebapi/Cilent"
	"wechatwebapi/Cilent/mm"
	"wechatwebapi/comm"
)

type UploadHeadImageParam struct {
	Wxid   string
	Base64 string
}

func UploadHeadImage(Data UploadHeadImageParam) wxClient.ResponseResult {
	var err error
	var protobufData []byte
	var errType int64

	D, err := comm.GetLoginData(Data.Wxid)
	if err != nil {
		return wxClient.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：%v", err.Error()),
			Data:    nil,
		}
	}

	ImgData := strings.Split(Data.Base64, ",")

	var ImgBase64 []byte

	if len(ImgData) > 1 {
		ImgBase64, _ = base64.StdEncoding.DecodeString(ImgData[1])
	} else {
		ImgBase64, _ = base64.StdEncoding.DecodeString(Data.Base64)
	}

	ImgStream := bytes.NewBuffer(ImgBase64)

	Startpos := 0
	datalen := 30000
	datatotalength := ImgStream.Len()

	ImgHash := wxClient.GetFileMD5Hash(ImgBase64)

	I := 0

	for {
		Startpos = I * datalen
		count := 0
		if datatotalength-Startpos > datalen {
			count = datalen
		} else {
			count = datatotalength - Startpos
		}
		if count < 0 {
			break
		}

		Databuff := make([]byte, count)
		_, _ = ImgStream.Read(Databuff)

		req := &mm.UploadHDHeadImgRequest{
			BaseRequest: &mm.BaseRequest{
				SessionKey:    D.SessionKey,
				Uin:           proto.Uint32(D.Uin),
				DeviceId:      D.Deviceid_byte,
				ClientVersion: proto.Int32(int32(wxClient.WxClientVersion)),
				DeviceType:    wxClient.DeviceTypeByte,
				Scene:         proto.Uint32(0),
			},
			TotalLen:    proto.Uint32(uint32(datatotalength)),
			StartPos:    proto.Uint32(uint32(Startpos)),
			HeadImgType: proto.Uint32(1),
			Data: &mm.SKBuiltinBufferT{
				ILen:   proto.Uint32(uint32(len(Databuff))),
				Buffer: Databuff,
			},
			ImgHash: proto.String(ImgHash),
		}

		//序列化
		reqData, _ := proto.Marshal(req)

		//发包
		protobufData, _, errType, err = comm.SendRequest(comm.SendPostData{
			Ip:            D.Mmtlsip,
			Cgiurl:        "/cgi-bin/micromsg-bin/uploadhdheadimg",
			Proxy:         D.Proxy,
			Encryption:    5,
			TwelveEncData: wxClient.PackSpecialCgiData{},
			PackData: wxClient.PackData{
				Reqdata:          reqData,
				Cgi:              157,
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
			break
		}

		I++
	}

	if err != nil {
		return wxClient.ResponseResult{
			Code:    errType,
			Success: false,
			Message: err.Error(),
			Data:    nil,
		}
	}

	//解包
	Response := mm.UploadHDHeadImgResponse{}
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
		Success: false,
		Message: "成功",
		Data:    &Response,
	}
}
