package wechat

type InvoiceUserInfo struct {
	Fee             uint32 `json:"fee"`
	Title           string `json:"title"`
	BillingTime     uint32 `json:"billing_time"`
	BillingNo       string `json:"billing_no"`
	BillingCode     string `json:"billing_code"`
	FeeWithoutTax   uint32 `json:"fee_without_tax"`
	Tax             uint32 `json:"tax"`
	PdfURL          string `json:"pdf_url"`
	ReimburseStatus string `json:"reimburse_status"`
	CheckCode       string `json:"check_code"`
}

type Invoice struct {
	CardID    string          `json:"card_id"`
	BeginTime uint32          `json:"begin_time"`
	EndTime   uint32          `json:"end_time"`
	OpenID    string          `json:"openid"`
	Type      string          `json:"type"`
	Payee     string          `json:"payee"`
	Detail    string          `json:"detail"`
	UserInfo  InvoiceUserInfo `json:"user_info"`
}

func (c *Client) GetInvoiceInfo(cardID, encryptCode string) (*Invoice, error) {
	const url = "https://api.weixin.qq.com/card/invoice/reimburse/getinvoiceinfo?access_token=ACCESS_TOKEN"
	var req struct {
		CardID      string `json:"card_id"`
		EncryptCode string `json:"encrypt_code"`
	}
	req.CardID = cardID
	req.EncryptCode = encryptCode

	var resp Invoice
	e := c.post(url, &req, &resp)
	return &resp, e
}
