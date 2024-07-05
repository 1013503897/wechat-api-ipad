package Group

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	wxClient "wechatwebapi/Cilent"
	"wechatwebapi/Cilent/mm"
	"wechatwebapi/comm"
)

type ScanIntoGroupParam struct {
	Wxid string
	Url  string
}

func ScanIntoGroup(Data ScanIntoGroupParam) wxClient.ResponseResult {
	D, err := comm.GetLoginata(Data.Wxid)
	if err != nil {
		return wxClient.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：%v", err.Error()),
			Data:    nil,
		}
	}

	//组包
	req := &mm.GetA8KeyReq{
		BaseRequest: &mm.BaseRequest{
			SessionKey:    D.Sessionkey,
			Uin:           proto.Uint32(D.Uin),
			DeviceId:      D.Deviceid_byte,
			ClientVersion: proto.Int32(int32(wxClient.WxClientVersion)),
			DeviceType:    wxClient.DeviceTypeByte,
			Scene:         proto.Uint32(0),
		},
		OpCode: proto.Uint32(2),
		ReqUrl: &mm.SKBuiltinStringT{
			String_: proto.String(Data.Url),
		},
		Scene:       proto.Uint32(4),
		FontScale:   proto.Uint32(100),
		NetType:     proto.String("WIFI"),
		CodeType:    proto.Uint32(19),
		CodeVersion: proto.Uint32(5),
		SubScene:    proto.Uint32(0),
	}

	//序列化
	reqData, _ := proto.Marshal(req)

	//发包
	protobufData, _, errType, err := comm.SendRequest(comm.SendPostData{
		Ip:            D.Mmtlsip,
		Cgiurl:        "/cgi-bin/micromsg-bin/geta8key",
		Proxy:         D.Proxy,
		Encryption:    5,
		TwelveEncData: wxClient.PackSpecialCgiData{},
		PackData: wxClient.PackData{
			Reqdata:          reqData,
			Cgi:              233,
			Uin:              D.Uin,
			Cookie:           D.Cooike,
			Sessionkey:       D.Sessionkey,
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
	GetA8KeyResp := mm.GetA8KeyResp{}
	err = proto.Unmarshal(protobufData, &GetA8KeyResp)
	if err != nil {
		return wxClient.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("反序列化失败：%v", err.Error()),
			Data:    nil,
		}
	}

	_, err = ScanIntoGrouppost(*GetA8KeyResp.FullURL)

	if strings.Index(err.Error(), "@chatroom") != -1 {
		return wxClient.ResponseResult{
			Code:    0,
			Success: true,
			Message: "进群成功",
			Data:    err.Error(),
		}
	}

	return wxClient.ResponseResult{
		Code:    -8,
		Success: false,
		Message: "进群失败",
		Data:    nil,
	}

}

func ScanIntoGrouppost(URL string) (string, error) {

	var err error

	postValue := url.Values{
		"forBlackberry": {"forceToUsePost"},
	}
	req, err := http.PostForm(URL, postValue)

	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Referer", URL)
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.116 Safari/537.36 QBCore/3.53.1159.400 QQBrowser/9.0.2524.400 Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.95 Safari/537.36 MicroMessenger/6.5.2.501 NetType/WIFI WindowsWechat")
	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
