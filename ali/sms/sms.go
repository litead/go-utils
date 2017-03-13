package sms

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync/atomic"
	"time"
)

type pair struct {
	key   string
	value string
}

type sender struct {
	appKey    string
	appSecret string
}

var (
	gNonce uint64
)

func NewSmsSender(appKey, appSecret string) *sender {
	return &sender{appKey: appKey, appSecret: appSecret}
}

func (s *sender) Send(mobile, signName, template, param string) error {
	params := []pair{
		{key: "AccessKeyId", value: s.appKey},
		{key: "Action", value: "SingleSendSms"},
		{key: "Format", value: "JSON"},
		{key: "ParamString", value: param},
		{key: "RecNum", value: mobile},
		{key: "SignName", value: signName},
		{key: "SignatureMethod", value: "HMAC-SHA1"},
		{key: "SignatureNonce", value: fmt.Sprintf("%d", getNonce())},
		{key: "SignatureVersion", value: "1.0"},
		{key: "TemplateCode", value: template},
		{key: "Timestamp", value: time.Now().UTC().Format("2006-01-02T15:04:05Z")},
		{key: "Version", value: "2016-09-27"},
	}

	buf := new(bytes.Buffer)
	for _, p := range params {
		buf.WriteString(url.QueryEscape(p.key))
		buf.WriteByte('=')
		buf.WriteString(url.QueryEscape(p.value))
		buf.WriteByte('&')
	}

	sign := getSignature(buf.Bytes()[:buf.Len()-1], s.appSecret)
	buf.Reset()
	buf.WriteString("Signature=")
	buf.WriteString(sign)
	for _, p := range params {
		buf.WriteByte('&')
		buf.WriteString(p.key)
		buf.WriteByte('=')
		buf.WriteString(p.value)
	}

	return doPost(buf)
}

func doPost(body *bytes.Buffer) error {
	const url = "http://sms.aliyuncs.com/"
	const bodyType = "application/x-www-form-urlencoded"

	resp, e := http.Post(url, bodyType, body)
	if e != nil {
		return e
	}

	var result struct {
		Model     string `json:"Model"`
		Message   string `json:"Message"`
		RequestID string `json:"RequestId"`
		HostID    string `json:"HostId"`
		Code      string `json:"Code"`
	}
	if e := json.NewDecoder(resp.Body).Decode(&result); e != nil {
		return e
	}

	if len(result.Code) != 0 {
		return fmt.Errorf("code: %s, message: %s", result.Code, result.Message)
	}

	return nil
}

func getSignature(data []byte, secret string) string {
	mac := hmac.New(sha1.New, []byte(secret+"&"))
	mac.Write([]byte("POST&%2F&"))
	str := url.QueryEscape(string(data))
	mac.Write([]byte(str))
	sig := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return url.QueryEscape(sig)
}

func getNonce() uint64 {
	nonce := atomic.AddUint64(&gNonce, 1)
	nonce = ((nonce & 0xf0f0f0f0f0f0f0f0) >> 4) | ((nonce & 0x0f0f0f0f0f0f0f0f) << 4)
	nonce = ((nonce & 0xffffffff00000000) >> 32) | ((nonce & 0x00000000ffffffff) << 32)
	return nonce
}

func init() {
	gNonce = uint64(time.Now().UnixNano())
}
