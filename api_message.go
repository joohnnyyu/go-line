package line

import (
	"context"
	"net/http"
)

type MessageService struct {
	client *Client
}

func (b *MessageService) Push(ctx context.Context, opt MessagePushOptions, options ...RequestOptionFunc) (*SentMessagesResponse, *Response, error) {
	req, err := b.client.NewRequest(ctx, http.MethodPost, "bot/message/push", opt, options)
	if err != nil {
		return nil, nil, err
	}

	m := new(SentMessagesResponse)
	resp, err := b.client.Do(req, m)
	if err != nil {
		return nil, nil, err
	}

	return m, resp, nil
}
