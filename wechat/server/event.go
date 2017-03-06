package server

type eventCommon struct {
	messageCommon
	Event string `xml:"Event"`
}

type SubscribeEvent struct {
	eventCommon
	EventKey string `xml:"EventKey"`
	Ticket   string `xml:"Ticket"`
}

type ScanEvent SubscribeEvent

type LocationEvent struct {
	eventCommon
	Latitude  float64 `xml:"Latitude"`
	Longitude float64 `xml:"Longitude"`
	Precision float64 `xml:"Precision"`
}

type ClickEvent SubscribeEvent
