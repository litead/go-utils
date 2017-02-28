package dayu

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
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

func NewSmsSender(appKey, appSecret string) *sender {
	return &sender{appKey: appKey, appSecret: appSecret}
}

func doPost(body *bytes.Buffer) error {
	const url = "http://gw.api.taobao.com/router/rest"
	const bodyType = "application/x-www-form-urlencoded"

	resp, e := http.Post(url, bodyType, body)
	if e != nil {
		return e
	}

	m := make(map[string]interface{})
	if e := json.NewDecoder(resp.Body).Decode(&m); e != nil {
		return e
	}

	if _, ok := m["alibaba_aliqin_fc_sms_num_send_response"]; ok {
		return nil
	}

	if m, ok := m["error_response"].(map[string]interface{}); ok {
		if s, ok := m["sub_msg"].(string); ok {
			return errors.New(s)
		}
	}

	return errors.New("发短信时发生未知错误")
}

func (s *sender) Send(mobile, signName, template, param string) error {
	params := []pair{
		{key: "app_key", value: s.appKey},
		{key: "format", value: "json"},
		{key: "method", value: "alibaba.aliqin.fc.sms.num.send"},
		{key: "rec_num", value: mobile},
		{key: "sign_method", value: "md5"},
		{key: "sms_free_sign_name", value: signName},
		{key: "sms_param", value: param},
		{key: "sms_template_code", value: template},
		{key: "sms_type", value: "normal"},
		{key: "timestamp", value: time.Now().Format("2006-01-02 15:04:05")},
		{key: "v", value: "2.0"},
	}

	buf := new(bytes.Buffer)
	buf.WriteString(s.appSecret)
	for _, p := range params {
		buf.WriteString(p.key)
		buf.WriteString(p.value)
	}
	buf.WriteString(s.appSecret)

	sign := fmt.Sprintf("%X", md5.Sum(buf.Bytes()))

	buf.Reset()
	buf.WriteString("sign=")
	buf.WriteString(sign)
	for _, p := range params {
		buf.WriteByte('&')
		buf.WriteString(url.QueryEscape(p.key))
		buf.WriteByte('=')
		buf.WriteString(url.QueryEscape(p.value))
	}

	return doPost(buf)
}
