package wechat

func (c *Client) CreateShortUrl(longUrl string) (string, error) {
	const url = "https://api.weixin.qq.com/cgi-bin/shorturl?access_token={ACCESS_TOKEN}"

	var req struct {
		Action  string `json:"action"`
		LongUrl string `json:"long_url"`
	}
	req.Action = "long2short"
	req.LongUrl = longUrl

	var resp struct {
		ShortUrl string `json:"short_url"`
	}
	e := c.post(url, &req, &resp)
	return resp.ShortUrl, e
}
