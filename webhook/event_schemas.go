package webhook

type FollowEvent struct {
	ReplyToken      string          `json:"replyToken,omitempty"`
	Type            string          `json:"type,omitempty"`
	Mode            string          `json:"mode,omitempty"`
	Timestamp       int64           `json:"timestamp,omitempty"`
	Source          Source          `json:"source,omitempty"`
	WebhookEventID  string          `json:"webhookEventId,omitempty"`
	DeliveryContext DeliveryContext `json:"deliveryContext,omitempty"`
	Follow          Follow          `json:"follow,omitempty"`
}

type UnFollowEvent struct {
	Type            string          `json:"type,omitempty"`
	Mode            string          `json:"mode,omitempty"`
	Timestamp       int64           `json:"timestamp,omitempty"`
	Source          Source          `json:"source,omitempty"`
	WebhookEventID  string          `json:"webhookEventId,omitempty"`
	DeliveryContext DeliveryContext `json:"deliveryContext,omitempty"`
}

type Source struct {
	Type   string `json:"type,omitempty"`
	UserID string `json:"userId,omitempty"`
}

type DeliveryContext struct {
	IsRedelivery bool `json:"isRedelivery,omitempty"`
}

type Follow struct {
	IsUnblocked bool `json:"isUnblocked,omitempty"`
}
