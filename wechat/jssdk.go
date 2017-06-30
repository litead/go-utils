package wechat

import (
	"crypto/sha1"
	"fmt"
)

type JSSDKConfig struct {
	Debug     bool     `json:"debug"`
	AppID     string   `json:"appId"`
	Timestamp uint32   `json:"timestamp"`
	Nonce     string   `json:"nonceStr"`
	Signature string   `json:"signature"`
	ApiList   []string `json:"jsApiList"`
}

func (c *Client) NewJSSDKConfig(timestamp uint32, nonce, url string) *JSSDKConfig {
	str := fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%v&url=%s",
		c.JSTicket.Value(),
		nonce,
		timestamp,
		url)
	data := sha1.Sum([]byte(str))

	return &JSSDKConfig{
		AppID:     c.AppID,
		Timestamp: timestamp,
		Nonce:     nonce,
		Signature: fmt.Sprintf("%02x", data),
	}
}
