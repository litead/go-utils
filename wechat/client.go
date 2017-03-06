package wechat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Client struct {
	AppID       string
	Secret      string
	lock        sync.RWMutex
	accessToken string
	expireTime  time.Time
}

func (c *Client) post(url string, req, resp interface{}) error {
	api = strings.Replace(url, "{ACCESS_TOKEN}", c.AccessToken())
	return c.doPost(api, req, resp)
}

func (c *Client) doPost(url string, req, resp interface{}) error {
	data, e := json.Marshal(req)
	if e != nil {
		return e
	}

	hresp, e := http.DefaultClient.Post(url, "application/json", bytes.NewReader(data))
	if e != nil {
		return e
	}

	return parseResponse(hresp, resp)
}

func (c *Client) get(url string, resp interface{}) error {
	api = strings.Replace(url, "{ACCESS_TOKEN}", c.AccessToken())
	return c.doGet(url, resp)
}

func (c *Client) doGet(url string, resp interface{}) error {
	hresp, e := http.DefaultClient.Get(url)
	if e != nil {
		return e
	}
	return parseResponse(hresp, resp)
}

func parseResponse(resp *http.Response, v interface{}) error {
	data, e := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if e != nil {
		return e
	}

	// first, is this an error?
	we := &Error{}
	if e := json.Unmarshal(data, we); e == nil && we.Code != 0 {
		return we
	}

	// no need to parse
	if v == nil {
		return nil
	}

	return json.Unmarshal(data, v)
}

func (c *Client) updateAccessToken() error {
	const urlFmt = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"

	url := fmt.Sprintf(urlFmt, c.AppID, c.Secret)
	var r struct {
		Token     string `json:"access_token"`
		ExpiresIn uint32 `json:"expires_in"`
	}
	if e := c.doGet(url, &r); e != nil {
		return e
	}

	// deduct 30 seconds, so that the token will be refreshed a little earlier
	// and never really expirs
	if r.ExpiresIn > 100 {
		r.ExpiresIn -= 30
	}

	c.accessToken = r.Token
	c.expireTime = time.Now().Add(time.Duration(r.ExpiresIn) * time.Second)
}

func (c *Client) UpdateAccessToken(force bool) (string, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if (!force) && time.Now().Before(c.ExpireTime) {
		return c.accessToken, nil
	}

	e := c.updateAccessToken()
	return c.accessToken, e
}

func (c *Client) SetAccessToken(accessToken string, expireTime time.Time) {
	c.lock.Lock()
	c.accessToken = accessToken
	c.expireTime = expireTime
	c.lock.Unlock()
}

func (c *Client) AccessToken() string {
	now := time.Now()

	c.lock.RLock()
	token, expire := c.accessToken, c.expireTime
	c.lock.RUnlock()

	if expire.After(now) {
		return token
	}

	c.lock.Lock()
	defer c.lock.Unlock()
	if c.expireTime.Before(now) {
		c.updateAccessToken()
	}

	return c.accessToken
}

func NewClient(appID, secret string) (*Client, error) {
	c := &Client{AppID: appID, Secret: secret}
	_, e := c.UpdateAccessToken(true)
	return c, e
}
