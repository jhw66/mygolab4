package webchat

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jhw66/myvideo_lab4/pkg/db/model"
	"github.com/jhw66/myvideo_lab4/pkg/db/repository/chat"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
	sendBufferSize = 256
)

type Client struct {
	hub       *Hub
	roomID    string
	userID    string
	nickname  string
	conn      *websocket.Conn
	send      chan []byte
	closeOnce sync.Once
}

func NewClient(hub *Hub, conn *websocket.Conn, roomID, userID, nickname string) *Client {
	return &Client{
		hub:      hub,
		roomID:   roomID,
		userID:   userID,
		nickname: nickname,
		conn:     conn,
		send:     make(chan []byte, sendBufferSize),
	}
}

func (c *Client) ServeClient(ctx context.Context, conn *websocket.Conn, chatRepo chat.ChatRepository) {
	select {
	case c.hub.register <- c:
	case <-c.hub.done:
		conn.Close()
		return
	}

	historyMessages, err := chatRepo.MessagesHistory(ctx, c.roomID, 100)
	if err != nil {
		log.Println("get history messages failed:", err)
		return
	}
	for _, message := range historyMessages {
		playload, err := json.Marshal(message)
		if err != nil {
			log.Println("marshal message failed:", err)
			continue
		}
		select {
		case c.send <- playload:
		default:
			log.Println("send channel is full")
			continue
		}
	}

	go c.writePump()
	c.readPump()
}

func (c *Client) closeSend() {
	c.closeOnce.Do(func() {
		close(c.send)
	})
}

func (c *Client) readPump() {
	defer func() {
		select {
		case c.hub.unregister <- c:
		default:
			log.Println("unregister client failed")
		}
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		var inmessage InMessage
		err := c.conn.ReadJSON(&inmessage)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		if inmessage.Type != MessageTypeChat {
			continue
		}
		cotent := strings.TrimSpace(inmessage.Content)
		if cotent == "" {
			continue
		}
		c.hub.toInbound(model.ChatMessage{
			RoomID:  c.roomID,
			UserID:  c.userID,
			Content: cotent,
		})
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				log.Println("send channel closed")
				c.conn.WriteMessage(websocket.CloseMessage, nil)
				return
			}
			err := c.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Println("write message failed:", err)
				continue
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println("write ping message failed:", err)
				return
			}
		}
	}
}
