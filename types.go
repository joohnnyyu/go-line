package line

type UserProfile struct {
	DisplayName   string `json:"displayName,omitempty"`
	UserID        string `json:"userId,omitempty"`
	PictureURL    string `json:"pictureUrl,omitempty"`
	StatusMessage string `json:"statusMessage,omitempty"`
	Language      string `json:"language,omitempty"`
}

type SentMessage struct {
	ID         string `json:"id"`
	QuoteToken string `json:"quoteToken"`
}

type SentMessagesResponse struct {
	SentMessages []SentMessage `json:"sentMessages"`
}

type MessagePushOptions struct {
	To       string    `json:"to,omitempty"`
	Messages []Message `json:"messages,omitempty"`
}

type Message struct {
	Type string `json:"type,omitempty"`
	Text string `json:"text,omitempty"`
}
