package wechat

import "fmt"

type Error struct {
	Code    int32  `json:"errcode"`
	Message string `json:"errmsg"`
}

func (e Error) Error() string {
	return fmt.Sprintf("[wechat] code: %d, message: %s", e.Code, e.Message)
}
