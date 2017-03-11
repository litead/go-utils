package wechat

import (
	"bytes"
	"encoding/xml"
	"regexp"
	"time"
)

///////////////////////////////////////////////////////////////////////////////
// common message bodies

type Media struct {
	MediaID string `xml:"MediaID" json:"media_id"`
}

type Video struct {
	Media
	Title        string `xml:"Title" json:"title"`
	Description  string `xml:"Description" json:"description"`
	ThumbMediaID string `xml:"-" json:"thumb_media_id"`
}

type Music struct {
	Title        string `xml:"Title" json:"title"`
	URL          string `xml:"MusicUrl" json:"musicurl"`
	HQURL        string `xml:"HQMusicUrl" json:"hqmusicurl"`
	Description  string `xml:"Description" json:"description"`
	ThumbMediaID string `xml:"ThumbMediaId" json:"thumb_media_id"`
}

type Article struct {
	Title       string `xml:"Title" json:"title"`
	Description string `xml:"Description" json:"description"`
	ImageURL    string `xml:"PicUrl" json:"picurl"`
	URL         string `xml:"Url" json:"url"`
}

///////////////////////////////////////////////////////////////////////////////
// base of all messages

type MessageBase struct {
	ToUserName   string `xml:"ToUserName" json:"touser"`
	FromUserName string `xml:"FromUserName" json:"-"`
	CreateTime   uint32 `xml:"CreateTime" json:"-"`
	MsgType      string `xml:"MsgType" json:"msgtype"`
}

func (m *MessageBase) init(to, msgType string) {
	m.ToUserName = to
	m.MsgType = msgType
}

func (m *MessageBase) initPassive(from, to, msgType string) {
	m.init(to, msgType)
	m.FromUserName = from
	m.CreateTime = uint32(time.Now().Unix())
}

///////////////////////////////////////////////////////////////////////////////
// passive messages

type PassiveMessageBase struct {
	XMLName struct{} `xml:"xml"`
	MessageBase
	MsgID uint64 `xml:"MsgId"`
}

type PassiveTextMessage struct {
	PassiveMessageBase
	Content string `xml:"Content"`
}

type PassiveImageMessage struct {
	PassiveMessageBase
	MediaID  string `xml:"MediaId"`
	ImageURL string `xml:"PicUrl"`
}

type PassiveVoiceMessage struct {
	PassiveMessageBase
	MediaID     string `xml:"MediaId"`
	Format      string `xml:"Format"`
	Recognition string `xml:"Recognition"` // when voice recognition enabled
}

type PassiveVideoMessage struct {
	PassiveMessageBase
	MediaID      string `xml:"MediaId"`
	ThumbMediaID string `xml:"ThumbMediaId"`
}
type PassiveShortVideoMessage PassiveVideoMessage

type PassiveLocationMessage struct {
	PassiveMessageBase
	Latitude  float64 `xml:"Location_X"`
	Longitude float64 `xml:"Location_Y"`
	Scale     float64 `xml:"Scale"`
	Label     string  `xml:"Label"`
}

type PassiveLinkMessage struct {
	PassiveMessageBase
	URL         string `xml:"Url"`
	Description string `xml:"Description"`
}

///////////////////////////////////////////////////////////////////////////////
// events

type EventBase struct {
	XMLName struct{} `xml:"xml"`
	MessageBase
	Event string `xml:"Event"`
}

type SubscribeEvent struct {
	EventBase
	EventKey string `xml:"EventKey"`
	Ticket   string `xml:"Ticket"`
}

type ScanEvent SubscribeEvent

type LocationEvent struct {
	EventBase
	Latitude  float64 `xml:"Latitude"`
	Longitude float64 `xml:"Longitude"`
	Precision float64 `xml:"Precision"`
}

type ClickEvent SubscribeEvent

type ViewEvent struct {
	SubscribeEvent
	MenuID string `xml:"MenuID"`
}

///////////////////////////////////////////////////////////////////////////////
// passive message replies

type PassiveReplyBase struct {
	XMLName struct{} `xml:"xml"`
	MessageBase
}

type PassiveTextReply struct {
	PassiveReplyBase
	Content string `xml:"Content"`
}

func NewPassiveTextReply(from, to string) *PassiveTextReply {
	r := &PassiveTextReply{}
	r.initPassive(from, to, "text")
	return r
}

type PassiveImageReply struct {
	PassiveReplyBase
	Image Media `xml:"Image"`
}

func NewPassiveImageReply(from, to string) *PassiveImageReply {
	r := &PassiveImageReply{}
	r.initPassive(from, to, "image")
	return r
}

type PassiveVoiceReply struct {
	PassiveReplyBase
	Voice Media `xml:"Voice"`
}

func NewPassiveVoiceReply(from, to string) *PassiveVoiceReply {
	r := &PassiveVoiceReply{}
	r.initPassive(from, to, "voice")
	return r
}

type PassiveVideoReply struct {
	PassiveReplyBase
	Video Video `xml:"Video"`
}

func NewPassiveVideoReply(from, to string) *PassiveVideoReply {
	r := &PassiveVideoReply{}
	r.initPassive(from, to, "video")
	return r
}

type PassiveMusicReply struct {
	PassiveReplyBase
	Music Music `xml:"Music"`
}

func NewPassiveMusicReply(from, to string) *PassiveMusicReply {
	r := &PassiveMusicReply{}
	r.initPassive(from, to, "music")
	return r
}

type PassiveNewsReply struct {
	PassiveReplyBase
	ArticleCount uint32    `xml:"ArticleCount"`
	Articles     []Article `xml:"Articles>item"`
}

func NewPassiveNewsReply(from, to string) *PassiveNewsReply {
	r := &PassiveNewsReply{}
	r.initPassive(from, to, "news")
	return r
}

///////////////////////////////////////////////////////////////////////////////
// custom messages

type TextMessage struct {
	MessageBase
	Text struct {
		Content string `json:"content"`
	} `json:"text"`
}

func NewTextMessage(to string) *TextMessage {
	msg := &TextMessage{}
	msg.init(to, "text")
	return msg
}

type ImageMessage struct {
	MessageBase
	Image Media `json:"image"`
}

func NewImageMessage(to string) *ImageMessage {
	msg := &ImageMessage{}
	msg.init(to, "image")
	return msg
}

type VoiceMessage struct {
	MessageBase
	Voice Media `json:"voice"`
}

func NewVoiceMessage(to string) *VoiceMessage {
	msg := &VoiceMessage{}
	msg.init(to, "voice")
	return msg
}

type VideoMessage struct {
	MessageBase
	Video Video `json:"video"`
}

func NewVideoMessage(to string) *VideoMessage {
	msg := &VideoMessage{}
	msg.init(to, "video")
	return msg
}

type MusicMessage struct {
	MessageBase
	Music Music `json:"music"`
}

func NewMusicMessage(to string) *MusicMessage {
	msg := &MusicMessage{}
	msg.init(to, "music")
	return msg
}

type NewsMessage struct {
	MessageBase
	News struct {
		Articles []Article `json:"articles"`
	} `json:"news"`
}

func NewNewsMessage(to string) *NewsMessage {
	msg := &NewsMessage{}
	msg.init(to, "news")
	return msg
}

type MPNewsMessage struct {
	MessageBase
	MPNews Media `json:"mpnews"`
}

func NewMPNewsMessage(to string) *MPNewsMessage {
	msg := &MPNewsMessage{}
	msg.init(to, "mpnews")
	return msg
}

type WXCardMessage struct {
	MessageBase
	Card struct {
		ID string `json:"card_id"`
	} `json:"wxcard"`
}

func NewWXCardMessage(to string) *WXCardMessage {
	msg := &WXCardMessage{}
	msg.init(to, "wxcard")
	return msg
}

///////////////////////////////////////////////////////////////////////////////

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

func (c *Client) SendCustomMessage(msg interface{}) error {
	const url = "https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token=ACCESS_TOKEN"
	return c.post(url, msg, nil)
}
