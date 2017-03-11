package wechat

import (
	"crypto/sha1"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

type Error struct {
	Code    int32  `json:"errcode"`
	Message string `json:"errmsg"`
}

func (e Error) Error() string {
	return fmt.Sprintf("[wechat] code: %d, message: %s", e.Code, e.Message)
}

func VerifySignature(token, timestamp, nonce, signature string) bool {
	strs := []string{token, timestamp, nonce}
	sort.Strings(strs)
	sig1 := sha1.Sum([]byte(strings.Join(strs, "")))
	if sig2, e := hex.DecodeString(signature); e == nil {
		return subtle.ConstantTimeCompare(sig1[:], sig2) == 1
	}
	return false
}
