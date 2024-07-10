package comm

import (
	"encoding/binary"
	"errors"
	"wechatwebapi/Cilent"
	"wechatwebapi/Cilent/mm"
	"wechatwebapi/Mmtls"
	"wechatwebapi/models"
)

type SendPostData struct {
	Ip            string
	Host          string
	Cgiurl        string
	Proxy         models.ProxyInfo
	Encryption    int
	TwelveEncData Cilent.PackSpecialCgiData
	PackData      Cilent.PackData
}

func MmtlsInitialize(Proxy models.ProxyInfo) (*Mmtls.HttpClientModel, *Mmtls.MmtlsClient, error) {
	//生成mmtls的公私钥。
	V1PrivKey, V1pubKey := Cilent.GenECDH415Key()
	V2PrivKey, V2pubKey := Cilent.GenECDH415Key()

	Shakehandpubkey := Mmtls.Shakehandpubkey{
		V1PrivKey: V1PrivKey,
		V1pubKey:  V1pubKey,
		V2PrivKey: V2PrivKey,
		V2pubKey:  V2pubKey,
	}

	httpclient := Mmtls.GenNewHttpClient(nil)
	MmtlsClient, err := httpclient.InitMmtlsShake(Cilent.MmtlsIp, Cilent.MmtlsHost, Proxy, Shakehandpubkey)

	if err != nil {
		return nil, &Mmtls.MmtlsClient{}, err
	}

	return httpclient, MmtlsClient, nil
}

func SendRequest(SENP SendPostData, MmtlsClient *Mmtls.MmtlsClient) (protobufData, Cookie []byte, errType int64, err error) {

	/*ttpclient := GenNewHttpClient()

	mmtlsret, err := ttpclient.InitMmtlsShake(SENP.Ip, SENP.Host, SENP.Proxy, SENP.Shakehandpubkey)

	if !mmtlsret {
		return nil, nil, -8, err
	}*/

	/*defer func() {
		err := recover()
		if err != nil {
			log.Error(err)
		}
	}()*/

	var sendData []byte

	//组包加密方式
	if SENP.Encryption == 12 {
		sendData = Cilent.PackSpecialCgi(SENP.TwelveEncData)
	} else {
		sendData = Cilent.Pack(SENP.PackData.Reqdata, SENP.PackData.Cgi, SENP.PackData.Uin, SENP.PackData.SessionKey, SENP.PackData.Cookie, SENP.PackData.Clientsessionkey, SENP.PackData.Loginecdhkey, SENP.PackData.EncryptType, SENP.PackData.UseCompress)
	}

	if SENP.Ip == "" {
		SENP.Ip = Cilent.MmtlsIp
	}

	if SENP.Host == "" {
		SENP.Host = Cilent.MmtlsHost
	}

	httpclient := Mmtls.GenNewHttpClient(MmtlsClient)

	response, err := httpclient.MMtlsPost(SENP.Ip, SENP.Host, SENP.Cgiurl, sendData, SENP.Proxy)

	if err != nil {
		return nil, nil, -1, err
	}

	if len(response) > 31 {
		//数据包解密/解包方式
		if SENP.Cgiurl == "/cgi-bin/micromsg-bin/newsync" {
			protobufData = Cilent.UnpackBusinessPacketWithAesGcm(response, SENP.PackData.Uin, &Cookie, SENP.PackData.Serversessionkey)
		} else {
			if SENP.Encryption == 12 {
				protobufData = Cilent.UnpackBusinessHybridEcdhPacket(response, 0, &Cookie, SENP.TwelveEncData.HybridEcdhPrivkey)
			} else {
				protobufData = Cilent.UnpackBusinessPacket(response, SENP.PackData.SessionKey, SENP.PackData.Uin, &Cookie)
			}
		}

		if protobufData != nil {
			return
		} else {
			return nil, nil, -8, errors.New("数据解密失败")
		}

	}

	Ret, err := RetConst(response)

	if Ret == -13 {
		return nil, nil, Ret, errors.New("您已退出微信")
	}

	return nil, nil, Ret, errors.New("微信服务返回信息：" + err.Error())

}

func RetConst(data []byte) (int64, error) {
	var Ret int32
	Ret = BytesToInt32(data[2:10])
	return int64(Ret), errors.New(mm.RetConst_name[BytesToInt32(data[2:10])])
}

func BytesToInt32(buf []byte) int32 {
	return int32(binary.BigEndian.Uint32(buf))
}
