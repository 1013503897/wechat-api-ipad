package Mmtls

import (
	"bytes"
	"encoding/binary"
	"errors"
	"strconv"
	"time"
	"wechatwebapi/models"
)

// mmtlspost
func (httpclient *HttpClientModel) MMtlsPost(ip, host, cgiurl string, data []byte, P models.ProxyInfo) ([]byte, error) {
	var err error
	newSendData := new(bytes.Buffer)
	binary.Write(newSendData, binary.BigEndian, int16(len(cgiurl)))
	newSendData.Write([]byte(cgiurl))
	binary.Write(newSendData, binary.BigEndian, int16(len(host)))
	newSendData.Write([]byte(host))
	binary.Write(newSendData, binary.BigEndian, int32(len(data)))
	newSendData.Write(data)
	sendData := new(bytes.Buffer)
	binary.Write(sendData, binary.BigEndian, int32(newSendData.Len()))
	sendData.Write(newSendData.Bytes())
	encryptData := httpclient.MmtlsEncryptData(sendData.Bytes())
	if encryptData == nil {
		return []byte{}, errors.New("MMTLS: 数据[EncryptData]失败")
	}
	var recvData []byte

	uniquenumstr := "/mmtls/" + strconv.Itoa(int(time.Now().Unix()))

	recvData, err = httpclient.POST(ip, uniquenumstr, encryptData, host, P)

	if err != nil {
		return []byte{}, err
	}

	response := new(bytes.Buffer)
	/*Separate := Separate(recv_data)
	for _, v := range Separate {
		response.Write(httpclient.MmtlsDecryptData(v))
	}*/

	response.Write(httpclient.MmtlsDecryptData(recvData))

	if response.Bytes() == nil {
		return []byte{}, errors.New("MMTLS: 数据[DecryptData]失败")
	}

	return response.Bytes(), nil
}
