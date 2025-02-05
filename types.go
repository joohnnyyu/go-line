package line

type UserProfile struct {
	DisplayName   string `json:"displayName,omitempty"`
	UserID        string `json:"userId,omitempty"`
	PictureURL    string `json:"pictureUrl,omitempty"`
	StatusMessage string `json:"statusMessage,omitempty"`
	Language      string `json:"language,omitempty"`
}

type MessagesResponse struct {
	SentMessages []SentMessage `json:"sentMessages"`
}

type SentMessage struct {
	ID         string `json:"id"`
	QuoteToken string `json:"quoteToken"`
}

type ValidatePushResponse struct {
	Message string        `json:"message"` // 主错误消息
	Details []ErrorDetail `json:"details"` // 错误详情
}

type ErrorDetail struct {
	Message  string `json:"message"`  // 错误消息
	Property string `json:"property"` // 相关属性
}
