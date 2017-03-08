package wechat

import (
	"fmt"
)

type UserInfo struct {
	Subscribe     uint8    `json:"subscribe"`
	Gender        uint8    `json:"sex"`
	OpenID        string   `json:"openid"`
	Nickname      string   `json:"nickname"`
	Language      string   `json:"language"`
	City          string   `json:"city"`
	Province      string   `json:"province"`
	Country       string   `json:"country"`
	HeadImageURL  string   `json:"headimgurl"`
	SubscribeTime uint32   `json:"subscribe_time"`
	GroupID       uint32   `json:"groupid"`
	UnionID       string   `json:"unionid"`
	Remark        string   `json:"remark"`
	Privileges    []string `json:"privilege"`
}

type OAuth2AccessToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    uint32 `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"openid"`
	Scope        string `json:"scope"`
	UnionID      string `json:"unionid"`
}

func (c *Client) GetOAuth2AccessToken(code string) (*OAuth2AccessToken, error) {
	const urlFmt = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
	url := fmt.Sprintf(urlFmt, c.AppID, c.Secret, code)
	token := &OAuth2AccessToken{}
	if e := c.doGet(url, token); e != nil {
		return nil, e
	}
	return token, nil
}

func (c *Client) RefreshOAuth2AccessToken(refreshToken string) (*OAuth2AccessToken, error) {
	const urlFmt = "https://api.weixin.qq.com/sns/oauth2/refresh_token?appid=%s&grant_type=refresh_token&refresh_token=%s"
	url := fmt.Sprintf(urlFmt, c.AppID, refreshToken)
	token := &OAuth2AccessToken{}
	if e := c.doGet(url, token); e != nil {
		return nil, e
	}
	return token, nil
}

func (c *Client) GetUserInfoByOAuth2AccessToken(accessToken, openID string) (*UserInfo, error) {
	const urlFmt = "https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN"
	url := fmt.Sprintf(urlFmt, accessToken, openID)
	ui := &UserInfo{}
	if e := c.doGet(url, ui); e != nil {
		return nil, e
	}
	return ui, nil
}

func (c *Client) GetUserInfoViaOAuth2(code string) (*UserInfo, error) {
	t, e := c.GetOAuth2AccessToken(code)
	if e != nil {
		return nil, e
	}
	return c.GetUserInfoByOAuth2AccessToken(t.AccessToken, t.OpenID)
}
