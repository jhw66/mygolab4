package webchat

import "time"

const (
	MessageTypeChat = "message"
	OnlineTTL       = 10 * time.Minute
)

type InMessage struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}
