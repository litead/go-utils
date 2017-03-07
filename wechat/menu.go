package wechat

type SubmenuItem struct {
	Type    string `json:"type,omitempty"`
	Name    string `json:"name"`
	Key     string `json:"key,omitempty"`
	URL     string `json:"url,omitempty"`
	MediaID string `json:"media_id,omitempty"`
}

type MenuItem struct {
	SubmenuItem
	Submenu []*SubmenuItem `json:"sub_button,omitempty"`
}

type MenuCondition struct {
	TagID              string `json:"tag_id,omitempty"`
	Gender             string `json:"sex,omitempty"`
	Country            string `json:"country,omitempty"`
	Province           string `json:"province,omitempty"`
	City               string `json:"city,omitempty"`
	ClientPlatformType string `json:"client_platform_type,omitempty"`
	Language           string `json:"language,omitempty"`
}

type AllMenu struct {
	Menu *struct {
		Button []*MenuItem `json:"button,omitempty"`
		ID     string      `json:"menuid,omitempty"`
	} `json:"menu,omitempty"`
	ConditionalMenu []struct {
		Button    []*MenuItem    `json:"button,omitempty"`
		MatchRule *MenuCondition `json:"matchrule,omitempty"`
		ID        string         `json:"menuid,omitempty"`
	}
}

const (
	urlCreateMenu         = "https://api.weixin.qq.com/cgi-bin/menu/create?access_token=ACCESS_TOKEN"
	urlGetAllMenu         = "https://api.weixin.qq.com/cgi-bin/menu/get?access_token=ACCESS_TOKEN"
	urlAddConditionalMenu = "https://api.weixin.qq.com/cgi-bin/menu/addconditional?access_token=ACCESS_TOKEN"
)

func (c *Client) CreateMenu(menu []*MenuItem) error {
	m := map[string][]*MenuItem{"button": menu}
	return c.post(urlCreateMenu, m, nil)
}

func (c *Client) CreateMenuFromJson(json string) error {
	return c.post(urlCreateMenu, json, nil)
}

func (c *Client) GetAllMenu() (*AllMenu, error) {
	m := &AllMenu{}
	if e := c.get(urlGetAllMenu, &m); e != nil {
		return nil, e
	}
	return m, nil
}

func (c *Client) GetAllMenuJson() (string, error) {
	var m string
	if e := c.get(urlGetAllMenu, &m); e != nil {
		return "", e
	}
	return m, nil
}

func (c *Client) DeleteAllMenu() error {
	const url = "https://api.weixin.qq.com/cgi-bin/menu/delete?access_token=ACCESS_TOKEN"
	return c.get(url, nil)
}

type menuID struct {
	ID string `json:"menuid"`
}

func (c *Client) AddConditionalMenu(menu []*MenuItem, condition *MenuCondition) (string, error) {
	m := map[string]interface{}{
		"button":    menu,
		"matchrule": condition,
	}
	var resp menuID
	if e := c.post(urlAddConditionalMenu, m, &resp); e != nil {
		return "", e
	}
	return resp.ID, nil
}

func (c *Client) AddConditionalMenuFromJson(json string) (string, error) {
	var resp menuID
	if e := c.post(urlAddConditionalMenu, json, &resp); e != nil {
		return "", e
	}
	return resp.ID, nil
}

func (c *Client) DeleteConditionalMenu(id string) error {
	const url = "https://api.weixin.qq.com/cgi-bin/menu/delconditional?access_token=ACCESS_TOKEN"
	req := menuID{ID: id}
	return c.post(url, &req, nil)
}
