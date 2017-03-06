package server

type messageCommon struct {
	ToUserName   string `xml:"ToUserName"`
	FromUserName string `xml:"FromUserName"`
	CreateTime   uint32 `xml:"CreateTime"`
	MsgType      string `xml:"MsgType"`
}

type message struct {
	messageCommon
	MsgID uint64 `xml:"MsgId"`
}

type TextMessage struct {
	message
	Content string `xml:"Content"`
}

type ImageMessage struct {
	message
	MediaID  string `xml:"MediaId"`
	ImageURL string `xml:"PicUrl"`
}

type VoiceMessage struct {
	MediaID     string `xml:"MediaId"`
	Format      string `xml:"Format"`
	Recognition string `xml:"Recognition"` // when voice recognition enabled
}

type VideoMessage struct {
	MediaID      string `xml:"MediaId"`
	ThumbMediaID string `xml:"ThumbMediaId"`
}
type ShortVideoMessage VideoMessage

type LocationMessage struct {
	Latitude  float64 `xml:"Location_X"`
	Longitude float64 `xml:"Location_Y"`
	Scale     float64 `xml:"Scale"`
	Label     string  `xml:"Label"`
}

type LinkMessage struct {
	URL         string `xml:"Url"`
	Description string `xml:"Description"`
}

type TextReply struct {
	messageCommon
	Content string `xml:"Content"`
}

type ImageReply struct {
	messageCommon
	Image struct {
		MediaID string `xml:"MediaId"`
	} `xml:"Image"`
}

type VoiceReply struct {
	messageCommon
	Voice struct {
		MediaID string `xml:"MediaId"`
	} `xml:"Voice"`
}

type VideoReply struct {
	messageCommon
	Video struct {
		MediaID     string `xml:"MediaId"`
		Title       string `xml:"Title"`
		Description string `xml:"Description"`
	} `xml:"Video"`
}

type MusicReply struct {
	messageCommon
	Music struct {
		Title        string `xml:"Title"`
		Description  string `xml:"Description"`
		URL          string `xml:"MusicUrl"`
		HQURL        string `xml:"HQMusicUrl"`
		ThumbMediaID string `xml:"ThumbMediaId"`
	} `xml:"Music"`
}

type NewsReply struct {
	messageCommon
	ArticleCount uint32 `xml:"ArticleCount"`
	Articles     []struct {
		Title       string `xml:"Title"`
		Description string `xml:"Description"`
		ImageURL    string `xml:"PicUrl"`
		URL         string `xml:"Url"`
	} `xml:"Articles>item"`
}