package wechat

import (
	"time"
)

type PassiveMessageBase struct {
	XMLName      struct{} `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`
	FromUserName string   `xml:"FromUserName"`
	CreateTime   uint32   `xml:"CreateTime"`
	MsgType      string   `xml:"MsgType"`
}

func (m *PassiveMessageBase) init(from, to, msgType string) {
	m.FromUserName = from
	m.ToUserName = to
	m.MsgType = msgType
	m.CreateTime = uint32(time.Now().Unix())
}

///////////////////////////////////////////////////////////////////////////////
// Messages

type PassiveMessage struct {
	PassiveMessageBase
	MsgID uint64 `xml:"MsgId"`
}

type PassiveTextMessage struct {
	PassiveMessage
	Content string `xml:"Content"`
}

type PassiveImageMessage struct {
	PassiveMessage
	MediaID  string `xml:"MediaId"`
	ImageURL string `xml:"PicUrl"`
}

type PassiveVoiceMessage struct {
	PassiveMessage
	MediaID     string `xml:"MediaId"`
	Format      string `xml:"Format"`
	Recognition string `xml:"Recognition"` // when voice recognition enabled
}

type PassiveVideoMessage struct {
	PassiveMessage
	MediaID      string `xml:"MediaId"`
	ThumbMediaID string `xml:"ThumbMediaId"`
}
type PassiveShortVideoMessage PassiveVideoMessage

type PassiveLocationMessage struct {
	PassiveMessage
	Latitude  float64 `xml:"Location_X"`
	Longitude float64 `xml:"Location_Y"`
	Scale     float64 `xml:"Scale"`
	Label     string  `xml:"Label"`
}

type PassiveLinkMessage struct {
	PassiveMessage
	URL         string `xml:"Url"`
	Description string `xml:"Description"`
}

///////////////////////////////////////////////////////////////////////////////
// Events

type EventBase struct {
	PassiveMessageBase
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
// Replies

type PassiveTextReply struct {
	PassiveMessageBase
	Content string `xml:"Content"`
}

func NewPassiveTextReply(from, to string) *PassiveTextReply {
	r := &PassiveTextReply{}
	r.init(from, to, "text")
	return r
}

type PassiveImageReply struct {
	PassiveMessageBase
	Image struct {
		MediaID string `xml:"MediaId"`
	} `xml:"Image"`
}

func NewPassiveImageReply(from, to string) *PassiveImageReply {
	r := &PassiveImageReply{}
	r.init(from, to, "image")
	return r
}

type PassiveVoiceReply struct {
	PassiveMessageBase
	Voice struct {
		MediaID string `xml:"MediaId"`
	} `xml:"Voice"`
}

func NewPassiveVoiceReply(from, to string) *PassiveVoiceReply {
	r := &PassiveVoiceReply{}
	r.init(from, to, "voice")
	return r
}

type PassiveVideoReply struct {
	PassiveMessageBase
	Video struct {
		MediaID     string `xml:"MediaId"`
		Title       string `xml:"Title"`
		Description string `xml:"Description"`
	} `xml:"Video"`
}

func NewPassiveVideoReply(from, to string) *PassiveVideoReply {
	r := &PassiveVideoReply{}
	r.init(from, to, "video")
	return r
}

type PassiveMusicReply struct {
	PassiveMessageBase
	Music struct {
		Title        string `xml:"Title"`
		Description  string `xml:"Description"`
		URL          string `xml:"MusicUrl"`
		HQURL        string `xml:"HQMusicUrl"`
		ThumbMediaID string `xml:"ThumbMediaId"`
	} `xml:"Music"`
}

func NewPassiveMusicReply(from, to string) *PassiveMusicReply {
	r := &PassiveMusicReply{}
	r.init(from, to, "music")
	return r
}

type PassiveNewsReply struct {
	PassiveMessageBase
	ArticleCount uint32    `xml:"ArticleCount"`
	Articles     []Article `xml:"Articles>item"`
}

func NewPassiveNewsReply(from, to string) *PassiveNewsReply {
	r := &PassiveNewsReply{}
	r.init(from, to, "news")
	return r
}
