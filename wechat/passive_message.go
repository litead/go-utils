package wechat

type PassiveMessageBase struct {
	XMLName      struct{} `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`
	FromUserName string   `xml:"FromUserName"`
	CreateTime   uint32   `xml:"CreateTime"`
	MsgType      string   `xml:"MsgType"`
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

type Event struct {
	PassiveMessageBase
	Event string `xml:"Event"`
}

type SubscribeEvent struct {
	Event
	EventKey string `xml:"EventKey"`
	Ticket   string `xml:"Ticket"`
}

type ScanEvent SubscribeEvent

type LocationEvent struct {
	Event
	Latitude  float64 `xml:"Latitude"`
	Longitude float64 `xml:"Longitude"`
	Precision float64 `xml:"Precision"`
}

type ClickEvent SubscribeEvent

///////////////////////////////////////////////////////////////////////////////
// Replies

type PassiveTextReply struct {
	PassiveMessageBase
	Content string `xml:"Content"`
}

type PassiveImageReply struct {
	PassiveMessageBase
	Image struct {
		MediaID string `xml:"MediaId"`
	} `xml:"Image"`
}

type PassiveVoiceReply struct {
	PassiveMessageBase
	Voice struct {
		MediaID string `xml:"MediaId"`
	} `xml:"Voice"`
}

type PassiveVideoReply struct {
	PassiveMessageBase
	Video struct {
		MediaID     string `xml:"MediaId"`
		Title       string `xml:"Title"`
		Description string `xml:"Description"`
	} `xml:"Video"`
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

type PassiveNewsReply struct {
	PassiveMessageBase
	ArticleCount uint32 `xml:"ArticleCount"`
	Articles     []struct {
		Title       string `xml:"Title"`
		Description string `xml:"Description"`
		ImageURL    string `xml:"PicUrl"`
		URL         string `xml:"Url"`
	} `xml:"Articles>item"`
}
