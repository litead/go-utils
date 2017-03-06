package server

import (
	"bytes"
	"crypto/sha1"
	"crypto/subtle"
	"encoding/hex"
	"encoding/xml"
	"regexp"
	"sort"
	"strings"
)

type MessageHandler interface {
	ProcessMessage(msg interface{}) interface{}
}

type Server struct {
	token          string
	encodingAESKey string
	messageHandler MessageHandler
}

func NewServer(token, encodingAESKey string, messageHandler MessageHandler) *Server {
	return &Server{
		token:          token,
		encodingAESKey: encodingAESKey,
		messageHandler: messageHandler,
	}
}

func (s *Server) VerifySignature(timestamp, nonce, signature string) bool {
	asig, e := hex.DecodeString(signature)
	if e != nil {
		return false
	}

	strs := []string{s.token, timestamp, nonce}
	sort.Strings(strs)
	dsig := sha1.Sum([]byte(strings.Join(strs, "")))

	// consider time used in network, I don't think we need to use
	// 'ConstantTimeCompare', but let's try to learn its usage
	if subtle.ConstantTimeCompare(asig, dsig[:]) != 1 {
		return false
	}

	return true
}

var (
	reMsgType   = regexp.MustCompile(`(?i)<\s*MsgType\s*>\s*<\s*(!\[CDATA\[)?\s*(\w+)\s*(\]\])?\s*><\s*/MsgType\s*>`)
	reEventType = regexp.MustCompile(`(?i)<\s*Event\s*>\s*<\s*(!\[CDATA\[)?\s*(\w+)\s*(\]\])?\s*><\s*/Event\s*>`)
)

func (s *Server) ProcessRequest(data []byte) []byte {
	matches := reMsgType.FindSubmatch(data)
	if len(matches) != 4 {
		return nil
	}

	var msg interface{}
	switch string(bytes.ToLower(matches[2])) {
	case "text":
		msg = &TextMessage{}
	case "image":
		msg = &ImageMessage{}
	case "voice":
		msg = &VoiceMessage{}
	case "video":
		msg = &VideoMessage{}
	case "shortvideo":
		msg = &ShortVideoMessage{}
	case "location":
		msg = &LocationMessage{}
	case "link":
		msg = &LinkMessage{}
	case "event":
		matches = reEventType.FindSubmatch(data)
		if len(matches) != 4 {
			return nil
		}
		switch string(bytes.ToLower(matches[2])) {
		case "subscribe":
		}
	}

	if msg == nil {
		return nil
	}

	if e := xml.Unmarshal(data, msg); e != nil {
		return nil
	}

	r := s.messageHandler.ProcessMessage(msg)
	if r == nil {
		return nil
	}

	if str, ok := r.(string); ok {
		return []byte(str)
	}

	data, _ = xml.Marshal(r)
	return data
}
