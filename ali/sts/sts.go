package sts

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sync/atomic"
	"time"
)

var (
	gNonce uint64
)

// Credential is the result of AssumeRole
type Credential struct {
	AccessKeyId     string
	AccessKeySecret string
	SecurityToken   string
	Expiration      string
}

func getSignature(str, secret string) string {
	data := []byte("GET&%2F&" + url.QueryEscape(str))
	mac := hmac.New(sha1.New, []byte(secret+"&"))
	mac.Write(data)
	sig := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return url.QueryEscape(sig)
}

func getNonce() uint64 {
	nonce := atomic.AddUint64(&gNonce, 1)
	nonce = ((nonce & 0xf0f0f0f0f0f0f0f0) >> 4) | ((nonce & 0x0f0f0f0f0f0f0f0f) << 4)
	nonce = ((nonce & 0xffffffff00000000) >> 32) | ((nonce & 0x00000000ffffffff) << 32)
	return nonce
}

func buildQueryString(accessKeyId, roleArn, session, policy string) string {

	// must in the dictionary order of argument name
	const format = "AccessKeyId=%s" +
		"&Action=AssumeRole" +
		"&DurationSeconds=1800" +
		"&Format=json" +
		"&Policy=%s" +
		"&RoleArn=%s" +
		"&RoleSessionName=%s" +
		"&SignatureMethod=HMAC-SHA1" +
		"&SignatureNonce=%d" +
		"&SignatureVersion=1.0" +
		"&Timestamp=%s" +
		"&Version=2015-04-01"

	s := fmt.Sprintf(format,
		url.QueryEscape(accessKeyId),
		url.QueryEscape(policy),
		url.QueryEscape(roleArn),
		url.QueryEscape(session),
		getNonce(),
		url.QueryEscape(time.Now().UTC().Format(time.RFC3339)),
	)

	return s
}

// AssumeRole give temporary permission to access OSS storage
func AssumeRole(accessKeyId, accessKeySecret, roleArn, session, policy string) (*Credential, error) {

	qs := buildQueryString(accessKeyId, roleArn, session, policy)
	qs += "&Signature=" + getSignature(qs, accessKeySecret)

	resp, err := http.Get("https://sts.aliyuncs.com/?" + qs)
	if err != nil {
		return nil, err
	}

	obj := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&obj)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(obj["Message"].(string))
	}

	c := obj["Credentials"].(map[string]interface{})
	return &Credential{
		AccessKeyId:     c["AccessKeyId"].(string),
		AccessKeySecret: c["AccessKeySecret"].(string),
		SecurityToken:   c["SecurityToken"].(string),
		Expiration:      c["Expiration"].(string),
	}, nil
}

func init() {
	gNonce = uint64(time.Now().UnixNano())
}
