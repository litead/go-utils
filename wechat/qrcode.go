package wechat

type QRcode struct {
	ExpireSeconds int    `json:"expire_seconds,omitempty"`
	Ticket        string `json:"ticket"`
	URL           string `json:"url"`
}

const urlCreateQRCode = "https://api.weixin.qq.com/cgi-bin/qrcode/create?access_token=ACCESS_TOKEN"

// CreateQRCode creates a QR code with the specified scene id.
// if expireSeconds is zero, create a permanent QR code,
// otherwise, create a temporary QR code with this expire seconds.
func (c *Client) CreateQRCode(sceneId uint32, expireSeconds uint32) (QRcode, error) {
	var req struct {
		ExpireSeconds uint32 `json:"expire_seconds,omitempty"`
		ActionName    string `json:"action_name"`
		ActionInfo    struct {
			Scene struct {
				SceneId uint32 `json:"scene_id"`
			} `json:"scene"`
		} `json:"action_info"`
	}

	if expireSeconds == 0 {
		req.ActionName = "QR_LIMIT_SCENE"
	} else {
		req.ActionName = "QR_SCENE"
		req.ExpireSeconds = expireSeconds
	}
	req.ActionInfo.Scene.SceneId = sceneId

	var qrcode QRcode
	e := c.post(urlCreateQRCode, &req, &qrcode)
	return qrcode, e
}

// CreateQRCodeStr creates a QR code with the specified scene id string.
// This function always create permanent QR code.
// Note the lenght of sceneId must between 1 and 64.
func (c *Client) CreateQRCodeStr(sceneId string) (QRcode, error) {
	var req struct {
		ActionName string `json:"action_name"`
		ActionInfo struct {
			Scene struct {
				SceneStr string `json:"scene_str"`
			} `json:"scene"`
		} `json:"action_info"`
	}

	req.ActionName = "QR_LIMIT_STR_SCENE"
	req.ActionInfo.Scene.SceneStr = sceneId

	var qrcode QRcode
	e := c.post(urlCreateQRCode, &req, &qrcode)
	return qrcode, e
}
