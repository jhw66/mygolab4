package webchat

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/jhw66/myvideo_lab4/pkg/cache/cacherepo"
	"github.com/jhw66/myvideo_lab4/pkg/db/model"
	"github.com/jhw66/myvideo_lab4/pkg/db/repository/chat"
)

type Hub struct {
	clientsByRoom map[string]map[*Client]bool
	register      chan *Client
	unregister    chan *Client

	inbound chan model.ChatMessage
	outbond chan model.ChatMessage

	done    chan struct{}
	stopped chan struct{}

	cancel map[string]func()
	once   sync.Once
}

func NewHub() *Hub {
	return &Hub{
		clientsByRoom: make(map[string]map[*Client]bool),
		register:      make(chan *Client),
		unregister:    make(chan *Client),
		inbound:       make(chan model.ChatMessage, 512),
		outbond:       make(chan model.ChatMessage, 512),
		done:          make(chan struct{}),
		stopped:       make(chan struct{}),
		cancel:        make(map[string]func()),
	}
}

func (h *Hub) Run(ctx context.Context, chatRepo chat.ChatRepository, chatCache cacherepo.ChatCache) {
	defer close(h.stopped)

	for {
		select {
		case client := <-h.register:
			roomClients := h.clientsByRoom[client.roomID]
			if roomClients == nil {
				roomClients = make(map[*Client]bool)
				h.clientsByRoom[client.roomID] = roomClients
				h.resigterRoom(ctx, client.roomID, chatCache)
			}
			roomClients[client] = true
			chatCache.MarkOnline(ctx, client.roomID, client.userID, OnlineTTL)
		case client := <-h.unregister:
			delete(h.clientsByRoom[client.roomID], client)
			client.closeSend()
			chatCache.MarkOffline(ctx, client.roomID, client.userID)
			if len(h.clientsByRoom[client.roomID]) == 0 {
				delete(h.clientsByRoom, client.roomID)
				h.unregisterRoom(client.roomID)
			}
		case message := <-h.inbound:
			if err := chatRepo.SaveMessage(ctx, &message); err != nil {
				log.Println("save message failed:", err)
				continue
			}
			if err := chatCache.PublishMessage(ctx, &message); err != nil {
				log.Println("publish message failed:", err)
				continue
			}
		case message := <-h.outbond:
			h.broadcast(ctx, message, chatCache)
		case <-h.done:
			h.closeAllClients(ctx, chatCache)
			h.cancelAllRooms()
			return
		}

	}
}

func (h *Hub) Stop() {
	close(h.done)
	<-h.stopped
}

func (h *Hub) toInbound(msg model.ChatMessage) {
	select {
	case h.inbound <- msg:
	default:
		log.Println("inbound channel is full")
	}
}

func (h *Hub) resigterRoom(ctx context.Context, roomID string, chatCache cacherepo.ChatCache) {
	if _, ok := h.cancel[roomID]; ok {
		return
	}
	ctx, cancel := context.WithCancel(ctx)
	stopSubscribe := h.subscribeRoom(ctx, roomID, chatCache)
	h.cancel[roomID] = func() {
		cancel()
		stopSubscribe()
	}
}

func (h *Hub) subscribeRoom(ctx context.Context, roomID string, chatCache cacherepo.ChatCache) func() {
	sub := chatCache.SubscribeMessages(ctx, roomID)
	done := make(chan struct{})

	go func() {
		defer close(done)
		ch := sub.Channel()
		for {
			select {
			case <-ctx.Done():
				return
			case rawMsg, ok := <-ch:
				if !ok {
					log.Println("subscribe room closed")
					return
				}
				var msg model.ChatMessage
				if err := json.Unmarshal([]byte(rawMsg.Payload), &msg); err != nil {
					log.Println("unmarshal message failed:", err)
					continue
				}
				h.toOutbound(msg)
			case <-h.stopped:
				return
			}
		}
	}()

	return func() {
		sub.Close()
		<-done
	}
}

func (h *Hub) toOutbound(msg model.ChatMessage) {
	select {
	case h.outbond <- msg:
	default:
		log.Println("outbound channel is full")
	}
}

func (h *Hub) unregisterRoom(roomID string) {
	if cancel, ok := h.cancel[roomID]; ok {
		cancel()
		delete(h.cancel, roomID)
	}

}

func (h *Hub) broadcast(ctx context.Context, msg model.ChatMessage, chatCache cacherepo.ChatCache) {
	payload, err := json.Marshal(msg)
	if err != nil {
		log.Println("marshal message failed:", err)
		return
	}

	roomClients := h.clientsByRoom[msg.RoomID]
	for client := range roomClients {
		chatCache.KeepOnline(ctx, msg.RoomID, client.userID, OnlineTTL)
		select {
		case client.send <- payload:
		default:
			log.Println("send channel is full")
		}
	}
}

func (h *Hub) closeAllClients(ctx context.Context, chatCache cacherepo.ChatCache) {
	for roomID, roomClients := range h.clientsByRoom {
		for client := range roomClients {
			client.closeSend()
			chatCache.MarkOffline(ctx, roomID, client.userID)
			delete(roomClients, client)
		}
		delete(h.clientsByRoom, roomID)
	}
}

func (h *Hub) cancelAllRooms() {
	cancellers := make([]func(), 0, len(h.cancel))
	for _, cancel := range h.cancel {
		cancellers = append(cancellers, cancel)
	}
	for _, cancel := range cancellers {
		cancel()
	}
	h.cancel = make(map[string]func())
}
