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
type errorResponse struct {
	ErrorResponse struct {
		Code    uint   `json:"code"`
		Msg     string `json:"msg"`
		SubCode string `json:"sub_code"`
		SumMsg  string `json:"sub_msg"`
	} `json:"error_response"`
}

func (e *errorResponse) Error() string {
	return fmt.Sprintf("[CODE]: %v, [MSG]: %v, [SUB_CODE]: %v, [SUB_MSG]: %v",
		e.ErrorResponse.Code,
		e.ErrorResponse.Msg,
		e.ErrorResponse.SubCode,
		e.ErrorResponse.SumMsg)
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

func NewClient(appKey, appSecret string, env uint) *Client {
	var baseURL string
	if env == EnvironmentSandBox {
		baseURL = "https://gw.api.tbsandbox.com/router/rest"
	} else if env == EnvironmentProduction {
		baseURL = "https://eco.taobao.com/router/rest"
	} else if env == EnvironmentOverseas {
		baseURL = "https://api.taobao.com/router/rest"
	} else {
		return nil
	}
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

func (c *Client) callAPI(method string, args []Argument, resp interface{}) error {
	args = append(args, Argument{Name: "method", Value: method},
		Argument{Name: "app_key", Value: c.appKey},
		Argument{Name: "format", Value: "json"},
		Argument{Name: "sign_method", Value: "hmac"},
		Argument{Name: "timestamp", Value: time.Now().Format("2006-01-02 15:04:05")},
		Argument{Name: "v", Value: "2.0"},
	)

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

	buf.Reset()
	for _, arg := range args {
		buf.WriteByte('&')
		buf.WriteString(url.QueryEscape(arg.Name))
		buf.WriteByte('=')
		buf.WriteString(url.QueryEscape(arg.Value))
	}
	data := buf.Bytes()
	data[0] = '?'

	hresp, e := http.DefaultClient.Get(c.baseURL + string(data))
	if e != nil {
		return e
	}

	data, e = ioutil.ReadAll(hresp.Body)
	hresp.Body.Close()
	if e != nil {
		return e
	}

	// first, is this an error?
	er := &errorResponse{}
	if e := json.Unmarshal(data, er); e != nil {
		return e
	} else if er.ErrorResponse.Code != 0 {
		return er
	}

	// no need to parse
	if resp == nil {
		return nil
	}

	return json.Unmarshal(data, resp)
}

type UatmTbkItem struct {
	ID          uint   `json:"num_iid"`
	Title       string `json:"title"`
	PictURL     string `json:"pict_url"`
	SmallImages struct {
		String []string `json:"string"`
	} `json:"small_images"`
	ReservePrice      float64 `json:"reserve_price,string"`
	ZKFinalPrice      float64 `json:"zk_final_price,string"`
	Provcity          string  `json:"provcity"`
	ItemURL           string  `json:"item_url"`
	ClickURL          string  `json:"click_url"`
	Nick              string  `json:"nick"`
	SellerID          uint    `json:"seller_id"`
	UserType          uint32  `json:"user_type"`
	Volume            uint32  `json:"volume"`
	TkRate            float64 `json:"tk_rate,string"`
	ZkFinalPriceWAP   float64 `json:"zk_final_price_wap,string"`
	ShopTitle         string  `json:"shop_title"`
	EventStartTime    string  `json:"event_start_time"`
	EventEndTime      string  `json:"event_end_time"`
	Type              uint32  `json:"type"`
	Status            uint32  `json:"status"`
	Category          uint    `json:"category"`
	CouponClickURL    string  `json:"coupon_click_url"`
	CouponEndTime     string  `json:"coupon_end_time"`
	CouponInfo        string  `json:"coupon_info"`
	CouponStartTime   string  `json:"coupon_start_time"`
	CouponTotalCount  uint32  `json:"coupon_total_count"`
	CouponRemainCount uint32  `json:"coupon_remain_count"`
}

func (c *Client) TBKGetUatmFavoritesItem(args []Argument) ([]UatmTbkItem, error) {
	var resp struct {
		Response struct {
			Results struct {
				Items []UatmTbkItem `json:"uatm_tbk_item"`
			} `json:"results"`
			TotalResults uint32 `json:"total_results"`
		} `json:"tbk_uatm_favorites_item_get_response"`
	}

	args = appendFieldsArgument(args, "num_iid,title,pict_url,small_images,reserve_price,zk_final_price,user_type,provcity,item_url,seller_id,volume,nick,shop_title,zk_final_price_wap,event_start_time,event_end_time,tk_rate,status,type")
	if e := c.callAPI("taobao.tbk.uatm.favorites.item.get", args, &resp); e != nil {
		return nil, e
	}

	return resp.Response.Results.Items, nil
}

type NTbkItem struct {
	ID          uint   `json:"num_iid"`
	Title       string `json:"title"`
	PictURL     string `json:"pict_url"`
	SmallImages struct {
		String []string `json:"string"`
	} `json:"small_images"`
	ReservePrice string `json:"reserve_price"`
	ZKFinalPrice string `json:"zk_final_price"`
	Provcity     string `json:"provcity"`
	ItemURL      string `json:"item_url"`
	Nick         string `json:"nick"`
	SellerID     uint   `json:"seller_id"`
	UserType     uint32 `json:"user_type"`
	Volume       uint32 `json:"volume"`
}

func (c *Client) TBKGetItem(args []Argument) ([]NTbkItem, error) {
	var resp struct {
		Response struct {
			Results struct {
				Items []NTbkItem `json:"n_tbk_item"`
			} `json:"results"`
			TotalResults uint32 `json:"total_results"`
		} `json:"tbk_item_get_response"`
	}

	args = appendFieldsArgument(args, "num_iid,title,pict_url,small_images,reserve_price,zk_final_price,user_type,provcity,item_url,seller_id,volume,nick")
	if e := c.callAPI("taobao.tbk.item.get", args, &resp); e != nil {
		return nil, e
	}

	return resp.Response.Results.Items, nil
}
