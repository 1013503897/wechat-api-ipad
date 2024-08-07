package Msg

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"time"
	wxClient "wechatwebapi/Cilent"
	"wechatwebapi/Cilent/mm"
	"wechatwebapi/comm"

	"github.com/golang/protobuf/proto"
)

type SendVideoMsgParam struct {
	Wxid      string
	ToWxid    string
	Video     string
	Image     string
	VideoTime uint32
}

func SendVideoMsg(Data SendVideoMsgParam) wxClient.ResponseResult {
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

	VideoBuffer, _ := base64.StdEncoding.DecodeString(Data.Video)
	ImageBuffer, _ := base64.StdEncoding.DecodeString(Data.Image)

	VideoStream := bytes.NewBuffer(VideoBuffer)
	ImageStream := bytes.NewBuffer(ImageBuffer)

	Startpos := 0
	datalen := 65000
	videoTotalLength := VideoStream.Len()
	imageTotalLength := ImageStream.Len()

	ClientImgId := fmt.Sprintf("%v_%v", Data.Wxid, time.Now().Unix())

	I := 0
	J := 0

	for {
		Startpos = I * datalen
		count := 0
		if videoTotalLength-Startpos > datalen {
			count = datalen
		} else {
			count = videoTotalLength - Startpos
		}
		if count < 0 {
			break
		}

		Databuff := make([]byte, count)
		_, _ = VideoStream.Read(Databuff)

		req := &mm.UploadVideoRequest{
			BaseRequest: &mm.BaseRequest{
				SessionKey:    D.SessionKey,
				Uin:           proto.Uint32(D.Uin),
				DeviceId:      D.Deviceid_byte,
				ClientVersion: proto.Int32(int32(wxClient.WxClientVersion)),
				DeviceType:    wxClient.DeviceTypeByte,
				Scene:         proto.Uint32(0),
			},
			ClientMsgID:  proto.String(ClientImgId),
			FromUserName: proto.String(Data.Wxid),
			ToUserName:   proto.String(Data.ToWxid),
			VideoData: &mm.SKBuiltinBufferT{
				ILen:   proto.Uint32(uint32(len(Databuff))),
				Buffer: Databuff,
			},
			VideoTotalLen: proto.Uint32(uint32(videoTotalLength)),
			VideoStartPos: proto.Uint32(uint32(Startpos)),
			ThumbData: &mm.SKBuiltinBufferT{
				ILen:   proto.Uint32(uint32(0)),
				Buffer: nil,
			},
			ThumbTotalLen: proto.Uint32(uint32(imageTotalLength)),
			ThumbStartPos: proto.Uint32(uint32(imageTotalLength)),
			PlayLength:    proto.Uint32(Data.VideoTime),
			NetWorkEnv:    proto.Uint32(uint32(1)),
			CameraType:    proto.Uint32(uint32(2)),
			FuncFlag:      proto.Uint32(uint32(2)),
			EncryVer:      proto.Int32(0),
			VideoFrom:     proto.Int32(0),
			ReqTime:       proto.Uint32(uint32(time.Now().Unix())),
		}

		//序列化
		reqData, _ := proto.Marshal(req)

		//发包
		protobufData, _, errType, err = comm.SendRequest(comm.SendPostData{
			Ip:            D.Mmtlsip,
			Cgiurl:        "/cgi-bin/micromsg-bin/uploadvideo",
			Proxy:         D.Proxy,
			Encryption:    5,
			TwelveEncData: wxClient.PackSpecialCgiData{},
			PackData: wxClient.PackData{
				Reqdata:          reqData,
				Cgi:              149,
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

	for {
		Startpos = J * datalen
		count := 0
		if imageTotalLength-Startpos > datalen {
			count = datalen
		} else {
			count = imageTotalLength - Startpos
		}
		if count < 0 {
			break
		}

		Databuff := make([]byte, count)
		_, _ = ImageStream.Read(Databuff)

		req := &mm.UploadVideoRequest{
			BaseRequest: &mm.BaseRequest{
				SessionKey:    D.SessionKey,
				Uin:           proto.Uint32(D.Uin),
				DeviceId:      D.Deviceid_byte,
				ClientVersion: proto.Int32(int32(wxClient.WxClientVersion)),
				DeviceType:    wxClient.DeviceTypeByte,
				Scene:         proto.Uint32(0),
			},
			ClientMsgID:  proto.String(ClientImgId),
			FromUserName: proto.String(Data.Wxid),
			ToUserName:   proto.String(Data.ToWxid),
			VideoData: &mm.SKBuiltinBufferT{
				ILen:   proto.Uint32(uint32(0)),
				Buffer: nil,
			},
			VideoTotalLen: proto.Uint32(uint32(videoTotalLength)),
			VideoStartPos: proto.Uint32(uint32(0)),
			ThumbData: &mm.SKBuiltinBufferT{
				ILen:   proto.Uint32(uint32(len(Databuff))),
				Buffer: Databuff,
			},
			ThumbTotalLen: proto.Uint32(uint32(Startpos)),
			ThumbStartPos: proto.Uint32(uint32(imageTotalLength)),
			PlayLength:    proto.Uint32(Data.VideoTime),
			NetWorkEnv:    proto.Uint32(uint32(1)),
			CameraType:    proto.Uint32(uint32(2)),
			FuncFlag:      proto.Uint32(uint32(2)),
			EncryVer:      proto.Int32(0),
			VideoFrom:     proto.Int32(0),
			ReqTime:       proto.Uint32(uint32(time.Now().Unix())),
		}

		//序列化
		reqData, _ := proto.Marshal(req)

		//发包
		protobufData, _, errType, err = comm.SendRequest(comm.SendPostData{
			Ip:            D.Mmtlsip,
			Cgiurl:        "/cgi-bin/micromsg-bin/uploadvideo",
			Proxy:         D.Proxy,
			Encryption:    5,
			TwelveEncData: wxClient.PackSpecialCgiData{},
			PackData: wxClient.PackData{
				Reqdata:          reqData,
				Cgi:              149,
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
	Response := mm.UploadVideoResponse{}
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
