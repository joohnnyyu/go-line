package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
)

// EventType represents the type of webhook event.
type EventType string

const (
	// ========================================
	// 需求/任务/缺陷类
	// ========================================

	EventTypeFollow   EventType = "follow"
	EventTypeUnFollow EventType = "unfollow"
)

func (e EventType) String() string {
	return string(e)
}

func validateWebhookEvent(secret string, signature string, body []byte) error {
	decoded, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return errors.New("line: webhook decoded error")
	}

	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write(body)

	if hmac.Equal(decoded, hash.Sum(nil)) {
		return nil
	}
	return errors.New("line: webhook compare error")
}

// ParseWebhookEvent parses the webhook event from the payload.
func ParseWebhookEvent(secret string, signature string, body []byte) ([]any, error) {
	err := validateWebhookEvent(secret, signature, body)
	if err != nil {
		return nil, err
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}

	// check event
	events, ok := raw["events"].([]interface{})
	if !ok || len(events) == 0 {
		return nil, errors.New("line: webhook events not found or empty")
	}
	var parsedEvents []any

	// get webhook type
	for _, event := range events {
		eventMap, ok := event.(map[string]interface{})
		if !ok {
			return nil, errors.New("line: event is not a valid map")
		}

		eventType, ok := eventMap["type"].(string)
		if !ok {
			return nil, errors.New("line: event type not found or not a string")
		}

		eventStr, err := json.Marshal(eventMap)
		if err != nil {
			return nil, err
		}
		switch EventType(eventType) {
		case EventTypeFollow:
			followEvent, err := decodeWebhookEvent[FollowEvent](eventStr)
			if err != nil {
				return nil, err
			}
			parsedEvents = append(parsedEvents, followEvent)
		case EventTypeUnFollow:
			unFollowEvent, err := decodeWebhookEvent[UnFollowEvent](eventStr)
			if err != nil {
				return nil, err
			}
			parsedEvents = append(parsedEvents, unFollowEvent)
		default: // todo: add more event types
			log.Printf("line: webhook event type [%s] not supported, skipping", eventType)
			continue
		}
	}
	return parsedEvents, nil
}

func decodeWebhookEvent[T any](body []byte) (*T, error) {
	var event T
	if err := json.Unmarshal(body, &event); err != nil {
		return nil, err
	}
	return &event, nil
}
