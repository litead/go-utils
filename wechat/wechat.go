package wechat

import (
	"bytes"
	"crypto/sha1"
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

func CalculateSignature(token, timestamp, nonce string) []byte {
	strs := []string{token, timestamp, nonce}
	sort.Strings(strs)
	sig := sha1.Sum([]byte(strings.Join(strs, "")))
	return sig[:]
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
		case "subscribe":
			msg = &SubscribeEvent{}
		case "scan":
			msg = &ScanEvent{}
		case "location":
			msg = &LocationEvent{}
		case "click":
			msg = &ClickEvent{}
		}
	}

	if msg != nil {
		if e := xml.Unmarshal(data, msg); e != nil {
			return nil
		}
	}

	return msg
}
