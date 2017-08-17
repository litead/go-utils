package wechat

import (
	"crypto/sha1"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type JSSDKConfig struct {
	Beta      bool   `json:"beta"`
	Debug     bool   `json:"debug"`
	AppID     string `json:"appId"`
	Timestamp uint32 `json:"timestamp"`
	Nonce     string `json:"nonceStr"`
	Signature string `json:"signature"`
	// exclude API list because this need to be configured in
	// web pages at most of the time
	//	ApiList   []string `json:"jsApiList"`
}

func (c *Client) NewJSSDKConfig(beta, debug bool, url string) *JSSDKConfig {
	cfg := JSSDKConfig{
		Beta:      beta,
		Debug:     debug,
		AppID:     c.AppID,
		Timestamp: uint32(time.Now().Unix()),
		Nonce:     strconv.FormatUint(rand.Uint64(), 36),
	}

	str := fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%v&url=%s",
		c.JSTicket.Value(),
		cfg.Nonce,
		cfg.Timestamp,
		url)
	cfg.Signature = fmt.Sprintf("%02x", sha1.Sum([]byte(str)))

	return &cfg
}
