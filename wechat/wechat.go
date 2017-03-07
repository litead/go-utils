package wechat

import (
	"bytes"
	"crypto/sha1"
	"crypto/subtle"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"regexp"
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

type Article struct {
	Title       string `xml:"Title" json:"title"`
	Description string `xml:"Description" json:"description"`
	ImageURL    string `xml:"PicUrl" json:"picurl"`
	URL         string `xml:"Url" json:"url"`
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

var (
	reMsgType   = regexp.MustCompile(`(?i)<\s*MsgType\s*>\s*<\s*(!\[CDATA\[)?\s*(\w+)\s*(\]\])?\s*><\s*/MsgType\s*>`)
	reEventType = regexp.MustCompile(`(?i)<\s*Event\s*>\s*<\s*(!\[CDATA\[)?\s*(\w+)\s*(\]\])?\s*><\s*/Event\s*>`)
)

func ParsePassiveMessage(data []byte) interface{} {
	matches := reMsgType.FindSubmatch(data)
	if len(matches) != 4 {
		return nil
	}

	var msg interface{}
	switch string(bytes.ToLower(matches[2])) {
	case "text":
		msg = &PassiveTextMessage{}
	case "image":
		msg = &PassiveImageMessage{}
	case "voice":
		msg = &PassiveVoiceMessage{}
	case "video":
		msg = &PassiveVideoMessage{}
	case "shortvideo":
		msg = &PassiveShortVideoMessage{}
	case "location":
		msg = &PassiveLocationMessage{}
	case "link":
		msg = &PassiveLinkMessage{}
	case "event":
		matches = reEventType.FindSubmatch(data)
		if len(matches) != 4 {
			return nil
		}
		switch string(bytes.ToLower(matches[2])) {
		case "subscribe", "unsubscribe":
			msg = &SubscribeEvent{}
		case "scan":
			msg = &ScanEvent{}
		case "location":
			msg = &LocationEvent{}
		case "click":
			msg = &ClickEvent{}
		case "view":
			msg = &ViewEvent{}
		}
	}

	if msg != nil {
		if e := xml.Unmarshal(data, msg); e != nil {
			return nil
		}
	}

	return msg
}
