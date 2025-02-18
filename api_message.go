package line

import (
	"context"
	"net/http"
)

type MessageService struct {
	client *Client
}

func (b *MessageService) Push(ctx context.Context, opt MessagePushOptions, options ...RequestOptionFunc) (*MessagesResponse, *Response, error) {
	req, err := b.client.NewRequest(ctx, http.MethodPost, "bot/message/push", opt, options)
	if err != nil {
		return nil, nil, err
	}

	m := new(MessagesResponse)
	resp, err := b.client.Do(req, m)
	if err != nil {
		return nil, nil, err
	}

	return m, resp, nil
}

func (b *MessageService) ValidatePush(ctx context.Context, opt ValidateMessagePushOptions, options ...RequestOptionFunc) (*ValidatePushResponse, *Response, error) {
	req, err := b.client.NewRequest(ctx, http.MethodPost, "bot/message/validate/push", opt, options)
	if err != nil {
		return nil, nil, err
	}

	m := new(ValidatePushResponse)
	resp, err := b.client.Do(req, m)
	if err != nil {
		return nil, nil, err
	}

	return m, resp, nil
}
