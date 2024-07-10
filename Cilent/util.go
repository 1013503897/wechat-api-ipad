package Cilent

import "hash"

// 0x17000841 708
// 0x17000C2B 712
// 0x1800322B 850

var WxClientVersion = 0x1800322B
var MmtlsIp = "szshort.weixin.qq.com"
var MmtlsHost = "extshort.weixin.qq.com"
var DeviceTypeByte = []byte("iPad iOS13.3.1")
var DeviceTypeStr = "iPad iOS13.3.1"

var HybridDecryptHash hash.Hash
var HybridServerpubhashFinal hash.Hash

type PackSpecialCgiData struct {
	Reqdata                    []byte
	Cgi                        int
	Encrypttype                int
	Extenddata                 []byte
	Uin                        uint32
	Cookies                    []byte
	ClientVersion              int
	HybridEcdhPrivkey          []byte
	HybridEcdhPubkey           []byte
	HybridEcdhInitServerPubKey []byte
}

type PackData struct {
	Reqdata          []byte
	Cgi              int
	Uin              uint32
	Cookie           []byte
	ClientVersion    int
	SessionKey       []byte
	EncryptType      uint8
	Loginecdhkey     []byte
	Clientsessionkey []byte
	Serversessionkey []byte
	UseCompress      bool
}

type ResponseResult struct {
	Code    int64
	Success bool
	Message string
	Data    interface{}
}
