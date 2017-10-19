package taobao

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

const (
	EnvironmentSandBox = iota
	EnvironmentProduction
	EnvironmentOverseas
)

/*
{
	"error_response": {
		"code":50,
		"msg":"Remote service error",
		"sub_code":"isv.invalid-parameter",
		"sub_msg":"非法参数"
	}
}
*/
type Error struct {
	Code    uint   `json:"code"`
	Msg     string `json:"msg"`
	SubCode string `json:"sub_code"`
	SumMsg  string `json:"sub_msg"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("[CODE]: %v, [MSG]: %v, [SUB_CODE]: %v, [SUB_MSG]: %v",
		e.Code,
		e.Msg,
		e.SubCode,
		e.SumMsg)
}

type Argument struct {
	Name  string
	Value string
}

type Client struct {
	appKey    string
	appSecret string
	baseURL   string
}

func NewClient(baseURL, appKey, appSecret string) *Client {
	return &Client{appKey: appKey, appSecret: appSecret, baseURL: baseURL}
}

func appendFieldsArgument(args []Argument, fields string) []Argument {
	for i := 0; i < len(args); i++ {
		if args[i].Name == "fields" {
			return args
		}
	}
	return append(args, Argument{Name: "fields", Value: fields})
}

func (c *Client) appendSignatureArgument(args []Argument) []Argument {
	sort.Slice(args, func(i, j int) bool {
		return strings.Compare(args[i].Name, args[j].Name) < 0
	})

	var buf bytes.Buffer
	for _, arg := range args {
		buf.WriteString(arg.Name)
		buf.WriteString(arg.Value)
	}

	mac := hmac.New(md5.New, []byte(c.appSecret))
	mac.Write(buf.Bytes())
	sig := mac.Sum(nil)
	args = append(args, Argument{
		Name:  "sign",
		Value: strings.ToUpper(hex.EncodeToString(sig)),
	})

	return args
}

func (c *Client) sendRequest(args []Argument) ([]byte, error) {
	var buf bytes.Buffer
	for _, arg := range args {
		buf.WriteString(url.QueryEscape(arg.Name))
		buf.WriteByte('=')
		buf.WriteString(url.QueryEscape(arg.Value))
		buf.WriteByte('&')
	}
	data := buf.Bytes()
	data = data[0 : len(data)-1]

	resp, e := http.DefaultClient.Post(c.baseURL, "application/x-www-form-urlencoded", bytes.NewReader(data))
	if e != nil {
		return nil, e
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func (c *Client) call(method string, args []Argument, resp interface{}) error {
	args = append(args, Argument{Name: "method", Value: method},
		Argument{Name: "app_key", Value: c.appKey},
		Argument{Name: "format", Value: "json"},
		Argument{Name: "sign_method", Value: "hmac"},
		Argument{Name: "timestamp", Value: time.Now().Format("2006-01-02 15:04:05")},
		Argument{Name: "v", Value: "2.0"},
	)
	args = c.appendSignatureArgument(args)

	data, e := c.sendRequest(args)
	if e != nil {
		return e
	}

	// first, is this an error?
	if bytes.HasPrefix(data, []byte(`{"error_response":`)) {
		err := &Error{}
		data = data[18 : len(data)-1]
		if e := json.Unmarshal(data, err); e != nil {
			return e
		} else if err.Code != 0 {
			return err
		}
	}

	// no need to parse
	if resp == nil {
		return nil
	}

	data = data[bytes.IndexByte(data, ':')+1 : len(data)-1]
	return json.Unmarshal(data, resp)
}
