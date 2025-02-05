package line

type MessagePushOptions struct {
	To       string    `json:"to,omitempty"`
	Messages []Message `json:"messages,omitempty"`
}

type ValidateMessagePushOptions struct {
	Messages []Message `json:"messages,omitempty"`
}

type MessageType string
type TemplateType string

const (
	TextMessageType     MessageType = "text"
	StickerMessageType  MessageType = "sticker"
	ImageMessageType    MessageType = "image"
	VideoMessageType    MessageType = "video"
	AudioMessageType    MessageType = "audio"
	LocationMessageType MessageType = "location"
	TemplateMessageType MessageType = "template"

	ButtonTemplateType        TemplateType = "buttons"
	ConfirmTemplateType       TemplateType = "confirm"
	CarouselTemplateType      TemplateType = "carousel"
	ImageCarouselTemplateType TemplateType = "image_carousel"
)

type Message interface{}

// TextMessage https://developers.line.biz/en/reference/messaging-api/#text-message
type TextMessage struct {
	Type       MessageType `json:"type,omitempty"`
	Text       string      `json:"text,omitempty"`
	QuoteToken string      `json:"quoteToken"`
}

// EmojiMessage https://developers.line.biz/en/reference/messaging-api/#text-message
type EmojiMessage struct {
	Type   MessageType `json:"type,omitempty"`
	Text   string      `json:"text,omitempty"`
	Emojis []Emoji     `json:"emojis,omitempty"`
}

type Emoji struct {
	Index     int    `json:"index,omitempty"`
	ProductID string `json:"productId,omitempty"`
	EmojiID   string `json:"emojiId,omitempty"`
}

// StickerMessage https://developers.line.biz/en/reference/messaging-api/#sticker-message
type StickerMessage struct {
	Type       MessageType `json:"type"`
	PackageID  string      `json:"packageId"`
	StickerID  string      `json:"stickerId"`
	QuoteToken string      `json:"quoteToken,omitempty"`
}

// ImageMessage https://developers.line.biz/en/reference/messaging-api/#image-message
type ImageMessage struct {
	Type               MessageType `json:"type"`
	OriginalContentURL string      `json:"originalContentUrl"`
	PreviewImageURL    string      `json:"previewImageUrl"`
}

// VideoMessage https://developers.line.biz/en/reference/messaging-api/#video-message
type VideoMessage struct {
	Type               MessageType `json:"type"`
	OriginalContentURL string      `json:"originalContentUrl"`
	PreviewImageURL    string      `json:"previewImageUrl"`
	TrackingID         string      `json:"trackingId"`
}

// AudioMessage https://developers.line.biz/en/reference/messaging-api/#audio-message
type AudioMessage struct {
	Type               MessageType `json:"type"`
	OriginalContentURL string      `json:"originalContentUrl"`
	Duration           int         `json:"duration"` // Duration is typically represented in milliseconds
}

// LocationMessage https://developers.line.biz/en/reference/messaging-api/#location-message
type LocationMessage struct {
	Type      MessageType `json:"type"`
	Title     string      `json:"title"`
	Address   string      `json:"address"`
	Latitude  float64     `json:"latitude"`
	Longitude float64     `json:"longitude"`
}

//ImagemapMessage https://developers.line.biz/en/reference/messaging-api/#imagemap-message

// TemplateMessage https://developers.line.biz/en/reference/messaging-api/#template-messages
type Template interface{}

type TemplateMessage struct {
	Type     MessageType `json:"type"`
	AltText  string      `json:"altText"`
	Template Template    `json:"template"`
}

type ButtonTemplate struct {
	Type                 string        `json:"type"`
	ThumbnailImageURL    string        `json:"thumbnailImageUrl"`
	ImageAspectRatio     string        `json:"imageAspectRatio"`
	ImageSize            string        `json:"imageSize"`
	ImageBackgroundColor string        `json:"imageBackgroundColor"`
	Title                string        `json:"title"`
	Text                 string        `json:"text"`
	DefaultAction        DefaultAction `json:"defaultAction"`
	Actions              []Action      `json:"actions"`
}

type ConfirmTemplate struct {
	Type    string   `json:"type"`
	Text    string   `json:"text"`
	Actions []Action `json:"actions"`
}

type CarouselTemplate struct {
	Type             string           `json:"type"`
	Columns          []CarouselColumn `json:"columns"`
	ImageAspectRatio string           `json:"imageAspectRatio"`
	ImageSize        string           `json:"imageSize"`
}

type CarouselColumn struct {
	ThumbnailImageURL    string        `json:"thumbnailImageUrl"`
	ImageBackgroundColor string        `json:"imageBackgroundColor"`
	Title                string        `json:"title"`
	Text                 string        `json:"text"`
	DefaultAction        DefaultAction `json:"defaultAction"`
	Actions              []Action      `json:"actions"`
}

type ImageCarouselTemplate struct {
	Type    string                `json:"type"`
	Columns []ImageCarouselColumn `json:"columns"`
}

type ImageCarouselColumn struct {
	ImageURL string `json:"imageUrl"`
	Action   Action `json:"action"`
}

type DefaultAction struct {
	Type  string `json:"type"`
	Label string `json:"label"`
	URI   string `json:"uri"`
}

type Action struct {
	Type  string `json:"type"`
	Label string `json:"label,omitempty"`
	Data  string `json:"data,omitempty"`
	URI   string `json:"uri,omitempty"`
	Text  string `json:"text,omitempty"`
}
