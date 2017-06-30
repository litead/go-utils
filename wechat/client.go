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

type Ticket struct {
	update     func() error
	lock       sync.Mutex
	value      string
	expireTime time.Time
}

func (t *Ticket) set(value string, expiresIn uint32) {
	// deduct 30 seconds, so that the token will be refreshed a little earlier
	// and never really expirs
	if expiresIn > 100 {
		expiresIn -= 30
	}

	t.value = value
	t.expireTime = time.Now().Add(time.Duration(expiresIn) * time.Second)
}

func (t *Ticket) Update(force bool) (string, error) {
	t.lock.Lock()
	defer t.lock.Unlock()

	if (!force) && time.Now().Before(t.expireTime) {
		return t.value, nil
	}

	e := t.update()
	return t.value, e
}

func (t *Ticket) Set(value string, expireTime time.Time) {
	t.lock.Lock()
	t.value = value
	t.expireTime = expireTime
	t.lock.Unlock()
}

func (t *Ticket) Value() string {
	now := time.Now()

	t.lock.Lock()

	if t.expireTime.Before(now) {
		t.update()
	}
	value := t.value

	t.lock.Unlock()
	return value
}

type Client struct {
	AppID       string
	Secret      string
	AccessToken Ticket
	JSTicket    Ticket
	CardTicket  Ticket
}

func (c *Client) post(url string, req, resp interface{}) error {
	url = strings.Replace(url, "ACCESS_TOKEN", c.AccessToken.Value(), -1)
	return c.doPost(url, req, resp)
}

func (c *Client) doPost(url string, req, resp interface{}) error {
	var data []byte

	if s, ok := req.(string); ok {
		data = []byte(s)
	} else if d, e := json.Marshal(req); e != nil {
		return e
	} else {
		data = d
	}

	hresp, e := http.DefaultClient.Post(url, "application/json", bytes.NewReader(data))
	if e != nil {
		return e
	}

	return parseResponse(hresp, resp)
}

func (c *Client) get(url string, resp interface{}) error {
	url = strings.Replace(url, "ACCESS_TOKEN", c.AccessToken.Value(), -1)
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

	if s, ok := v.(*string); ok {
		*s = string(data)
		return nil
	}
	return json.Unmarshal(data, v)
}

func NewClient(appID, secret string) *Client {
	c := &Client{AppID: appID, Secret: secret}

	c.AccessToken.update = func() error {
		const urlFmt = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"

		url := fmt.Sprintf(urlFmt, c.AppID, c.Secret)
		var r struct {
			Token     string `json:"access_token"`
			ExpiresIn uint32 `json:"expires_in"`
		}
		if e := c.doGet(url, &r); e != nil {
			return e
		}

		c.AccessToken.set(r.Token, r.ExpiresIn)
		return nil
	}

	update := func(t *Ticket, url string) func() error {
		return func() error {
			var r struct {
				Ticket    string `json:"ticket"`
				ExpiresIn uint32 `json:"expires_in"`
			}
			if e := c.doGet(url, &r); e != nil {
				return e
			}
			t.set(r.Ticket, r.ExpiresIn)
			return nil
		}
	}

	c.JSTicket.update = update(&c.JSTicket, "https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=ACCESS_TOKEN&type=jsapi")
	c.CardTicket.update = update(&c.CardTicket, "https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=ACCESS_TOKEN&type=wx_card")

	return c
}
