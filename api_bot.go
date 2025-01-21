package line

import (
	"context"
	"fmt"
	"net/http"
)

type BotService struct {
	client *Client
}

func (b *BotService) Profile(ctx context.Context, userID string) (*UserProfile, *Response, error) {
	u := fmt.Sprintf("bot/profile/%s", userID)
	req, err := b.client.NewRequest(ctx, http.MethodGet, u, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	m := new(UserProfile)
	resp, err := b.client.Do(req, m)
	if err != nil {
		return nil, nil, err
	}

	return m, resp, nil
}
