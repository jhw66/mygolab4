package cacherepo

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/jhw66/myvideo_lab4/pkg/db/model"
	"github.com/redis/go-redis/v9"
)

func roomChannel(roomID string) string {
	return fmt.Sprintf("chat:room:%s:channel", roomID)
}

func roomOnlineKey(roomID, userID string) string {
	return fmt.Sprintf("chat:room:%s:online:%s", roomID, userID)
}

type ChatCache interface {
	PublishMessage(ctx context.Context, msg *model.ChatMessage) error
	SubscribeMessages(ctx context.Context, roomID string) *redis.PubSub
	MarkOnline(ctx context.Context, roomID, userID string, ttl time.Duration)
	MarkOffline(ctx context.Context, roomID, userID string)
	KeepOnline(ctx context.Context, roomID, userID string, ttl time.Duration)
	Close() error
}

type chatCache struct {
	rdb *redis.Client
}

func NewChatCache(rdb *redis.Client) ChatCache {
	return &chatCache{rdb: rdb}
}

func (c *chatCache) PublishMessage(ctx context.Context, msg *model.ChatMessage) error {
	playload, err := json.Marshal(msg)
	if err != nil {
		log.Printf("failed to marshal message: %v", err)
		return err
	}
	return c.rdb.Publish(ctx, roomChannel(msg.RoomID), playload).Err()
}

func (c *chatCache) SubscribeMessages(ctx context.Context, roomID string) *redis.PubSub {
	return c.rdb.Subscribe(ctx, roomChannel(roomID))
}

func (c *chatCache) MarkOnline(ctx context.Context, roomID, userID string, ttl time.Duration) {
	if err := c.rdb.Set(ctx, roomOnlineKey(roomID, userID), "1", ttl).Err(); err != nil {
		log.Printf("mark offline failed: %v", err)
	}
}
func (c *chatCache) MarkOffline(ctx context.Context, roomID, userID string) {
	if err := c.rdb.Del(ctx, roomOnlineKey(roomID, userID)).Err(); err != nil {
		log.Printf("mark offline failed: %v", err)
	}
}
func (c *chatCache) KeepOnline(ctx context.Context, roomID, userID string, ttl time.Duration) {
	if err := c.rdb.Expire(ctx, roomOnlineKey(roomID, userID), ttl).Err(); err != nil {
		log.Printf("keep online failed: %v", err)
	}
}
func (c *chatCache) Close() error {
	return c.rdb.Close()
}
