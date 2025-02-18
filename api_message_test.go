package line

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Push(t *testing.T) {
	// 创建一个模拟的 HTTP 服务器
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 验证请求路径和方法
		if r.URL.Path != "/v2/bot/message/push" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("unexpected method: %s", r.Method)
		}

		// 模拟成功的响应
		response := MessagesResponse{}
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			return
		}
	}))
	defer ts.Close()

	// 创建 Client 并设置 baseURL 为模拟服务器的 URL
	client, err := NewClient("test-token")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	err = client.setBaseURL(ts.URL)
	if err != nil {
		return
	}

	// 创建 MessageService
	messageService := &MessageService{client: client}

	// 调用 Push 方法
	opt := MessagePushOptions{
		To: "U1234567890",
		Messages: []Message{
			TextMessage{Type: TextMessageType, Text: "Hello, World!"},
		},
	}
	_, _, err = messageService.Push(context.Background(), opt)
	if err != nil {
		t.Errorf("Push returned error: %v", err)
	}
}

func Test_ValidatePush(t *testing.T) {
	// 创建一个模拟的 HTTP 服务器
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 验证请求路径和方法
		if r.URL.Path != "/v2/bot/message/validate/push" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("unexpected method: %s", r.Method)
		}

		// 模拟成功的响应
		response := ValidatePushResponse{}
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			return
		}
	}))
	defer ts.Close()

	// 创建 Client 并设置 baseURL 为模拟服务器的 URL
	client, err := NewClient("test-token")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	err = client.setBaseURL(ts.URL)
	if err != nil {
		return
	}

	// 创建 MessageService
	messageService := &MessageService{client: client}

	// 调用 Push 方法
	opt := ValidateMessagePushOptions{
		Messages: []Message{
			TextMessage{Type: TextMessageType, Text: "Hello, World!"},
		},
	}
	_, _, err = messageService.ValidatePush(context.Background(), opt)
	if err != nil {
		t.Errorf("ValidatePush returned error: %v", err)
	}
}

func TestParseTextMessage(t *testing.T) {
	message := `{"to":"U1234567890","messages":[{"type":"text","text":"Hello, World!"}]}`
	expected := MessagePushOptions{
		To: "U1234567890",
		Messages: []Message{
			TextMessage{Type: TextMessageType, Text: "Hello, World!"},
		},
	}

	opt, err := parseMessage(context.Background(), message)

	assert.NoError(t, err)
	assert.Equal(t, expected.To, opt.To)
	assert.Equal(t, len(expected.Messages), len(opt.Messages))
	assert.Equal(t, expected.Messages[0], opt.Messages[0])
}

func TestParseEmojiMessage(t *testing.T) {
	message := `{"to":"U1234567890","messages":[{"type":"text","text":"$ LINE emoji $","emojis":[{"index":0,"productId":"5ac1bfd5040ab15980c9b435","emojiId":"001"}]}]}`
	expected := MessagePushOptions{
		To: "U1234567890",
		Messages: []Message{
			EmojiMessage{
				Type: TextMessageType,
				Text: "$ LINE emoji $",
				Emojis: []Emoji{
					{Index: 0, ProductID: "5ac1bfd5040ab15980c9b435", EmojiID: "001"},
				},
			},
		},
	}

	opt, err := parseMessage(context.Background(), message)

	assert.NoError(t, err)
	assert.Equal(t, expected.To, opt.To)
	assert.Equal(t, len(expected.Messages), len(opt.Messages))
	assert.Equal(t, expected.Messages[0], opt.Messages[0])
}

func TestParseStickerMessage(t *testing.T) {
	message := `{"to":"U1234567890","messages":[{"type":"sticker","packageId":"446","stickerId":"1988"}]}`
	expected := MessagePushOptions{
		To: "U1234567890",
		Messages: []Message{
			StickerMessage{Type: StickerMessageType, PackageID: "446", StickerID: "1988"},
		},
	}

	opt, err := parseMessage(context.Background(), message)

	assert.NoError(t, err)
	assert.Equal(t, expected.To, opt.To)
	assert.Equal(t, len(expected.Messages), len(opt.Messages))
	assert.Equal(t, expected.Messages[0], opt.Messages[0])
}

func TestParseImageMessage(t *testing.T) {
	message := `{"to":"U1234567890","messages":[{"type":"image","originalContentUrl":"https://example.com/original.jpg","previewImageUrl":"https://example.com/preview.jpg"}]}`
	expected := MessagePushOptions{
		To: "U1234567890",
		Messages: []Message{
			ImageMessage{Type: ImageMessageType, OriginalContentURL: "https://example.com/original.jpg", PreviewImageURL: "https://example.com/preview.jpg"},
		},
	}

	opt, err := parseMessage(context.Background(), message)

	assert.NoError(t, err)
	assert.Equal(t, expected.To, opt.To)
	assert.Equal(t, len(expected.Messages), len(opt.Messages))
	assert.Equal(t, expected.Messages[0], opt.Messages[0])
}

func TestParseVideoMessage(t *testing.T) {
	message := `{"to":"U1234567890","messages":[{"type":"video","originalContentUrl":"https://example.com/original.mp4","previewImageUrl":"https://example.com/preview.jpg","trackingId":"track-id"}]}`
	expected := MessagePushOptions{
		To: "U1234567890",
		Messages: []Message{
			VideoMessage{Type: VideoMessageType, OriginalContentURL: "https://example.com/original.mp4", PreviewImageURL: "https://example.com/preview.jpg", TrackingID: "track-id"},
		},
	}

	opt, err := parseMessage(context.Background(), message)

	assert.NoError(t, err)
	assert.Equal(t, expected.To, opt.To)
	assert.Equal(t, len(expected.Messages), len(opt.Messages))
	assert.Equal(t, expected.Messages[0], opt.Messages[0])
}

func TestParseAudioMessage(t *testing.T) {
	message := `{"to":"U1234567890","messages":[{"type":"audio","originalContentUrl":"https://example.com/original.m4a","duration":60000}]}`
	expected := MessagePushOptions{
		To: "U1234567890",
		Messages: []Message{
			AudioMessage{Type: AudioMessageType, OriginalContentURL: "https://example.com/original.m4a", Duration: 60000},
		},
	}

	opt, err := parseMessage(context.Background(), message)

	assert.NoError(t, err)
	assert.Equal(t, expected.To, opt.To)
	assert.Equal(t, len(expected.Messages), len(opt.Messages))
	assert.Equal(t, expected.Messages[0], opt.Messages[0])
}

func TestParseLocationMessage(t *testing.T) {
	message := `{"to":"U1234567890","messages":[{"type":"location","title":"my location","address":"1-3 Kioicho, Chiyoda-ku, Tokyo, 102-8282, Japan","latitude":35.67966,"longitude":139.73669}]}`

	expected := MessagePushOptions{
		To: "U1234567890",
		Messages: []Message{
			LocationMessage{Type: LocationMessageType, Title: "my location", Address: "1-3 Kioicho, Chiyoda-ku, Tokyo, 102-8282, Japan", Latitude: 35.67966, Longitude: 139.73669},
		},
	}

	opt, err := parseMessage(context.Background(), message)

	assert.NoError(t, err)
	assert.Equal(t, expected.To, opt.To)
	assert.Equal(t, len(expected.Messages), len(opt.Messages))
	assert.Equal(t, expected.Messages[0], opt.Messages[0])
}

func TestParseTemplateMessage(t *testing.T) {
	message := `{
		"to":"U1234567890",
		"messages":[
			{"type":"template","altText":"This is a buttons template","template":{"type":"buttons","thumbnailImageUrl":"https://example.com/bot/images/image.jpg","imageAspectRatio":"rectangle","imageSize":"cover","imageBackgroundColor":"#FFFFFF","title":"Menu","text":"Please select","defaultAction":{"type":"uri","label":"View detail","uri":"http://example.com/page/123"},"actions":[{"type":"postback","label":"Buy","data":"action=buy&itemid=123"},{"type":"postback","label":"Add to cart","data":"action=add&itemid=123"},{"type":"uri","label":"View detail","uri":"http://example.com/page/123"}]}},
			{"type":"template","altText":"this is a confirm template","template":{"type":"confirm","text":"Are you sure?","actions":[{"type":"message","label":"Yes","text":"yes"},{"type":"message","label":"No","text":"no"}]}},
			{"type":"template","altText":"this is a carousel template","template":{"type":"carousel","columns":[{"thumbnailImageUrl":"https://example.com/bot/images/item1.jpg","imageBackgroundColor":"#FFFFFF","title":"this is menu","text":"description","defaultAction":{"type":"uri","label":"View detail","uri":"http://example.com/page/123"},"actions":[{"type":"postback","label":"Buy","data":"action=buy&itemid=111"},{"type":"postback","label":"Add to cart","data":"action=add&itemid=111"},{"type":"uri","label":"View detail","uri":"http://example.com/page/111"}]},{"thumbnailImageUrl":"https://example.com/bot/images/item2.jpg","imageBackgroundColor":"#000000","title":"this is menu","text":"description","defaultAction":{"type":"uri","label":"View detail","uri":"http://example.com/page/222"},"actions":[{"type":"postback","label":"Buy","data":"action=buy&itemid=222"},{"type":"postback","label":"Add to cart","data":"action=add&itemid=222"},{"type":"uri","label":"View detail","uri":"http://example.com/page/222"}]}],"imageAspectRatio":"rectangle","imageSize":"cover"}},
			{"type":"template","altText":"this is a image carousel template","template":{"type":"image_carousel","columns":[{"imageUrl":"https://example.com/bot/images/item1.jpg","action":{"type":"postback","label":"Buy","data":"action=buy&itemid=111"}},{"imageUrl":"https://example.com/bot/images/item2.jpg","action":{"type":"message","label":"Yes","text":"yes"}},{"imageUrl":"https://example.com/bot/images/item3.jpg","action":{"type":"uri","label":"View detail","uri":"http://example.com/page/222"}}]}}
		]
	}`
	expected := MessagePushOptions{
		To: "U1234567890",
		Messages: []Message{
			TemplateMessage{
				Type:    "template",
				AltText: "This is a buttons template",
				Template: ButtonTemplate{
					Type:                 "buttons",
					ThumbnailImageURL:    "https://example.com/bot/images/image.jpg",
					ImageAspectRatio:     "rectangle",
					ImageSize:            "cover",
					ImageBackgroundColor: "#FFFFFF",
					Title:                "Menu",
					Text:                 "Please select",
					DefaultAction: &DefaultAction{
						Type:  "uri",
						Label: "View detail",
						URI:   "http://example.com/page/123",
					},
					Actions: []Action{
						{
							Type:  "postback",
							Label: "Buy",
							Data:  "action=buy&itemid=123",
						},
						{
							Type:  "postback",
							Label: "Add to cart",
							Data:  "action=add&itemid=123",
						},
						{
							Type:  "uri",
							Label: "View detail",
							URI:   "http://example.com/page/123",
						},
					},
				},
			},
			TemplateMessage{
				Type:    "template",
				AltText: "this is a confirm template",
				Template: ConfirmTemplate{
					Type: "confirm",
					Text: "Are you sure?",
					Actions: []Action{
						{
							Type:  "message",
							Label: "Yes",
							Text:  "yes",
						},
						{
							Type:  "message",
							Label: "No",
							Text:  "no",
						},
					},
				},
			},
			TemplateMessage{
				Type:    "template",
				AltText: "this is a carousel template",
				Template: CarouselTemplate{
					Type:             "carousel",
					ImageAspectRatio: "rectangle",
					ImageSize:        "cover",
					Columns: []CarouselColumn{
						{
							ThumbnailImageURL:    "https://example.com/bot/images/item1.jpg",
							ImageBackgroundColor: "#FFFFFF",
							Title:                "this is menu",
							Text:                 "description",
							DefaultAction: &DefaultAction{
								Type:  "uri",
								Label: "View detail",
								URI:   "http://example.com/page/123",
							},
							Actions: []Action{
								{
									Type:  "postback",
									Label: "Buy",
									Data:  "action=buy&itemid=111",
								},
								{
									Type:  "postback",
									Label: "Add to cart",
									Data:  "action=add&itemid=111",
								},
								{
									Type:  "uri",
									Label: "View detail",
									URI:   "http://example.com/page/111",
								},
							},
						},
						{
							ThumbnailImageURL:    "https://example.com/bot/images/item2.jpg",
							ImageBackgroundColor: "#000000",
							Title:                "this is menu",
							Text:                 "description",
							DefaultAction: &DefaultAction{
								Type:  "uri",
								Label: "View detail",
								URI:   "http://example.com/page/222",
							},
							Actions: []Action{
								{
									Type:  "postback",
									Label: "Buy",
									Data:  "action=buy&itemid=222",
								},
								{
									Type:  "postback",
									Label: "Add to cart",
									Data:  "action=add&itemid=222",
								},
								{
									Type:  "uri",
									Label: "View detail",
									URI:   "http://example.com/page/222",
								},
							},
						},
					},
				},
			},
			TemplateMessage{
				Type:    "template",
				AltText: "this is a image carousel template",
				Template: ImageCarouselTemplate{
					Type: "image_carousel",
					Columns: []ImageCarouselColumn{
						{
							ImageURL: "https://example.com/bot/images/item1.jpg",
							Action: Action{
								Type:  "postback",
								Label: "Buy",
								Data:  "action=buy&itemid=111",
							},
						},
						{
							ImageURL: "https://example.com/bot/images/item2.jpg",
							Action: Action{
								Type:  "message",
								Label: "Yes",
								Text:  "yes",
							},
						},
						{
							ImageURL: "https://example.com/bot/images/item3.jpg",
							Action: Action{
								Type:  "uri",
								Label: "View detail",
								URI:   "http://example.com/page/222",
							},
						},
					},
				},
			},
		},
	}

	opt, err := parseMessage(context.Background(), message)

	assert.NoError(t, err)
	assert.Equal(t, expected.To, opt.To)
	assert.Equal(t, len(expected.Messages), len(opt.Messages))
	assert.Equal(t, expected.Messages[0], opt.Messages[0])
}

func parseMessage(_ context.Context, message string) (*MessagePushOptions, error) {
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
