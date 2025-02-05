package line

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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

func (b *MessageService) ParseMessage(_ context.Context, message string) (*MessagePushOptions, error) {
	var options MessagePushOptions
	if err := json.Unmarshal([]byte(message), &options); err != nil {
		return nil, err
	}

	for i, msg := range options.Messages {
		var base map[string]interface{}
		b, err := json.Marshal(msg)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal interface message: %v", err)
		}
		if err := json.Unmarshal(b, &base); err != nil {
			return nil, fmt.Errorf("failed to unmarshal base message: %v", err)
		}

		if msgType, ok := base["type"].(string); ok {
			switch MessageType(msgType) {
			case TextMessageType:
				if _, hasEmojis := base["emojis"]; hasEmojis {
					var emojiMsg EmojiMessage
					if err := json.Unmarshal(b, &emojiMsg); err != nil {
						return nil, fmt.Errorf("failed to unmarshal EmojiMessage: %v", err)
					}
					options.Messages[i] = emojiMsg
				} else {
					var textMsg TextMessage
					if err := json.Unmarshal(b, &textMsg); err != nil {
						return nil, fmt.Errorf("failed to unmarshal TextMessage: %v", err)
					}
					options.Messages[i] = textMsg
				}

			case StickerMessageType:
				var stickerMsg StickerMessage
				if err := json.Unmarshal(b, &stickerMsg); err != nil {
					return nil, fmt.Errorf("failed to unmarshal StickerMessage: %v", err)
				}
				options.Messages[i] = stickerMsg

			case ImageMessageType:
				var imageMsg ImageMessage
				if err := json.Unmarshal(b, &imageMsg); err != nil {
					return nil, fmt.Errorf("failed to unmarshal ImageMessage: %v", err)
				}
				options.Messages[i] = imageMsg
			case VideoMessageType:
				var videoMsg VideoMessage
				if err := json.Unmarshal(b, &videoMsg); err != nil {
					return nil, fmt.Errorf("failed to unmarshal VideoMessage: %v", err)
				}
				options.Messages[i] = videoMsg
			case AudioMessageType:
				var audioMsg AudioMessage
				if err := json.Unmarshal(b, &audioMsg); err != nil {
					return nil, fmt.Errorf("failed to unmarshal AudioMessage: %v", err)
				}
				options.Messages[i] = audioMsg
			case LocationMessageType:
				var locationMsg LocationMessage
				if err := json.Unmarshal(b, &locationMsg); err != nil {
					return nil, fmt.Errorf("failed to unmarshal LocationMessage: %v", err)
				}
				options.Messages[i] = locationMsg
			case TemplateMessageType:
				var templateMsg TemplateMessage
				if err := json.Unmarshal(b, &templateMsg); err != nil {
					return nil, fmt.Errorf("failed to unmarshal TemplateMessage: %v", err)
				}

				// 基于Type使用不同的结构体解码
				var template map[string]interface{}
				t, err := json.Marshal(templateMsg.Template)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal interface template message: %v", err)
				}
				if err := json.Unmarshal(t, &template); err != nil {
					return nil, fmt.Errorf("failed to unmarshal base template message: %v", err)
				}

				if templateType, ok := template["type"].(string); ok {
					switch TemplateType(templateType) {
					case ButtonTemplateType:
						var tmpl ButtonTemplate
						if err := json.Unmarshal(t, &tmpl); err != nil {
							log.Fatalf("Error parsing ButtonTemplate: %v", err)
						}
						templateMsg.Template = tmpl
					case ConfirmTemplateType:
						var tmpl ConfirmTemplate
						if err := json.Unmarshal(t, &tmpl); err != nil {
							log.Fatalf("Error parsing ConfirmTemplate: %v", err)
						}
						templateMsg.Template = tmpl
					case CarouselTemplateType:
						var tmpl CarouselTemplate
						if err := json.Unmarshal(t, &tmpl); err != nil {
							log.Fatalf("Error parsing CarouselTemplate: %v", err)
						}
						templateMsg.Template = tmpl
					case ImageCarouselTemplateType:
						var tmpl ImageCarouselTemplate
						if err := json.Unmarshal(t, &tmpl); err != nil {
							log.Fatalf("Error parsing ImageCarouselTemplate: %v", err)
						}
						templateMsg.Template = tmpl
					default:
						fmt.Println("Unknown template type")
					}
				}
				options.Messages[i] = templateMsg
			// Add additional cases for other types as needed
			default:
				return nil, fmt.Errorf("unknown message type: %s", msgType)
			}
		} else {
			return nil, fmt.Errorf("message type is missing or not a string")
		}
	}

	log.Printf("MessagePushOptions: %+v", options)
	return &options, nil
}
