package taobao

import "strconv"

type NTbkItem struct {
	ID           uint        `json:"num_iid"`
	Title        string      `json:"title"`
	PictURL      string      `json:"pict_url"`
	SmallImages  StringSlice `json:"small_images"`
	ReservePrice string      `json:"reserve_price"`
	ZKFinalPrice string      `json:"zk_final_price"`
	Provcity     string      `json:"provcity"`
	ItemURL      string      `json:"item_url"`
	Nick         string      `json:"nick"`
	SellerID     uint        `json:"seller_id"`
	UserType     uint32      `json:"user_type"`
	Volume       uint32      `json:"volume"`
}

func (c *Client) TbkGetItem(args []Argument, count int) ([]NTbkItem, error) {
	const fields = "num_iid,title,pict_url,small_images,reserve_price,zk_final_price," +
		"user_type,provcity,item_url,seller_id,volume,nick"
	args = appendFieldsArgument(args, fields)

	if count <= 0 {
		count = 0x7fffffff
	}

	pageSize := count
	if count >= 100 {
		pageSize = 100
	}
	args = append(args, Argument{Name: "page_size", Value: strconv.Itoa(pageSize)})

	result := make([]NTbkItem, 0, 1024)
	for page := 1; len(result) < count; page++ {
		var resp struct {
			Results struct {
				Items []NTbkItem `json:"n_tbk_item"`
			} `json:"results"`
			TotalResults int `json:"total_results"`
		}

		if e := c.call("taobao.tbk.item.get", args, &resp); e != nil {
			return nil, e
		}

		result = append(result, resp.Results.Items...)
		if pageSize > len(resp.Results.Items) || len(result) >= resp.TotalResults {
			break
		}
	}

	return result, nil
}

type TbkFavorites struct {
	Type  uint32 `json:"type"`
	ID    uint   `json:"favorites_id"`
	Title string `json:"favorites_title"`
}

func (c *Client) TbkGetUatmFavorites(args []Argument, count int) ([]TbkFavorites, error) {
	args = appendFieldsArgument(args, "favorites_title,favorites_id,type")

	if count <= 0 {
		count = 0x7fffffff
	}

	pageSize := count
	if count >= 100 {
		pageSize = 100
	}
	args = append(args, Argument{Name: "page_size", Value: strconv.Itoa(pageSize)})

	result := make([]TbkFavorites, 0, 128)
	for page := 1; len(result) < count; page++ {
		var resp struct {
			Results struct {
				Items []TbkFavorites `json:"tbk_favorites"`
			} `json:"results"`
			TotalResults int `json:"total_results"`
		}

		targs := make([]Argument, len(args))
		copy(targs, args)
		targs = append(targs, Argument{Name: "page_no", Value: strconv.Itoa(page)})

		if e := c.call("taobao.tbk.uatm.favorites.get", args, &resp); e != nil {
			return nil, e
		}

		result = append(result, resp.Results.Items...)
		if pageSize > len(resp.Results.Items) || len(result) >= resp.TotalResults {
			break
		}
	}

	return result, nil
}

type UatmTbkItem struct {
	ID                uint        `json:"num_iid"`
	Title             string      `json:"title"`
	PictURL           string      `json:"pict_url"`
	SmallImages       StringSlice `json:"small_images"`
	ReservePrice      float64     `json:"reserve_price,string"`
	ZKFinalPrice      float64     `json:"zk_final_price,string"`
	Provcity          string      `json:"provcity"`
	ItemURL           string      `json:"item_url"`
	ClickURL          string      `json:"click_url"`
	Nick              string      `json:"nick"`
	SellerID          uint        `json:"seller_id"`
	UserType          uint32      `json:"user_type"`
	Volume            uint32      `json:"volume"`
	TkRate            float64     `json:"tk_rate,string"`
	ZkFinalPriceWAP   float64     `json:"zk_final_price_wap,string"`
	ShopTitle         string      `json:"shop_title"`
	EventStartTime    Time        `json:"event_start_time"`
	EventEndTime      Time        `json:"event_end_time"`
	Type              uint32      `json:"type"`
	Status            uint32      `json:"status"`
	Category          uint        `json:"category"`
	CouponStartTime   Time        `json:"coupon_start_time"`
	CouponEndTime     Time        `json:"coupon_end_time"`
	CouponClickURL    string      `json:"coupon_click_url"`
	CouponInfo        CouponInfo  `json:"coupon_info"`
	CouponTotalCount  uint32      `json:"coupon_total_count"`
	CouponRemainCount uint32      `json:"coupon_remain_count"`
}

func (c *Client) TbkGetUatmFavoritesItem(args []Argument, count int) ([]UatmTbkItem, error) {
	const fields = "num_iid,title,pict_url,small_images,reserve_price,zk_final_price," +
		"user_type,provcity,item_url,seller_id,volume,nick,shop_title,zk_final_price_wap," +
		"event_start_time,event_end_time,tk_rate,status,type," +
		"coupon_click_url,coupon_end_time,coupon_info,coupon_start_time,coupon_total_count,coupon_remain_count"
	args = appendFieldsArgument(args, fields)

	if count <= 0 {
		count = 0x7fffffff
	}

	pageSize := count
	if count >= 100 {
		pageSize = 100
	}
	args = append(args, Argument{Name: "page_size", Value: strconv.Itoa(pageSize)})

	result := make([]UatmTbkItem, 0, 1024)
	for page := 1; len(result) < count; page++ {
		var resp struct {
			Results struct {
				Items []UatmTbkItem `json:"uatm_tbk_item"`
			} `json:"results"`
			TotalResults int `json:"total_results"`
		}

		targs := make([]Argument, len(args))
		copy(targs, args)
		targs = append(targs, Argument{Name: "page_no", Value: strconv.Itoa(page)})

		if e := c.call("taobao.tbk.uatm.favorites.item.get", targs, &resp); e != nil {
			return nil, e
		}

		result = append(result, resp.Results.Items...)
		if pageSize > len(resp.Results.Items) || len(result) >= resp.TotalResults {
			break
		}
	}

	return result, nil
}
