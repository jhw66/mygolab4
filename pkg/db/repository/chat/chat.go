package chat

import (
	"context"

	"github.com/jhw66/myvideo_lab4/pkg/db/model"
	"gorm.io/gorm"
)

const (
	VisibilityPublic  = "public"
	VisibilityPrivate = "private"
	RoleMember        = "member"
	RoleOwner         = "owner"
	RoleAdmin         = "admin"
)

type ChatRepository interface {
	CreateRoom(ctx context.Context, room *model.ChatRoom) error
	CreateRoomWithTx(ctx context.Context, tx *gorm.DB, room *model.ChatRoom) error
	FindRoomByRoomName(ctx context.Context, name string) (*model.ChatRoom, error)
	FindRoomByRoomId(ctx context.Context, id string) (*model.ChatRoom, error)
	JoinRoom(ctx context.Context, roomID, userID, role string) error
	JoinRoomWithTx(ctx context.Context, tx *gorm.DB, roomID, userID, role string) error
	LeaveRoom(ctx context.Context, roomID, userID string) error
	DeleteRoom(ctx context.Context, roomID string) error
	ListVisibleRooms(ctx context.Context, userID string) ([]model.ChatRoom, error)
	ListMyRooms(ctx context.Context, userID string) ([]model.ChatRoom, error)
	FindMemberById(ctx context.Context, roomID, userID string) (*model.ChatMember, error)
	ListRoomMembers(ctx context.Context, roomID string) ([]model.ChatMember, error)
	ListRoomMessages(ctx context.Context, roomID string, page, pageSize int) ([]model.ChatMessage, error)
	CountRoomMessages(ctx context.Context, roomID string) (int64, error)
	MessagesHistory(ctx context.Context, roomID string, limit int) ([]model.ChatMessage, error)
	SaveMessage(ctx context.Context, message *model.ChatMessage) error
}

type chatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) ChatRepository {
	return &chatRepository{
		db: db,
	}
}

func (r *chatRepository) CreateRoom(ctx context.Context, room *model.ChatRoom) error {
	if err := r.db.WithContext(ctx).Create(room).Error; err != nil {
		return err
	}
	return nil
}

func (r *chatRepository) CreateRoomWithTx(ctx context.Context, tx *gorm.DB, room *model.ChatRoom) error {
	if err := tx.WithContext(ctx).Create(room).Error; err != nil {
		return err
	}
	return nil
}

func (r *chatRepository) FindRoomByRoomName(ctx context.Context, name string) (*model.ChatRoom, error) {
	var room model.ChatRoom
	if err := r.db.WithContext(ctx).Where("room_name = ?", name).Take(&room).Error; err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *chatRepository) FindRoomByRoomId(ctx context.Context, id string) (*model.ChatRoom, error) {
	var room model.ChatRoom
	if err := r.db.WithContext(ctx).Where("id = ?", id).Take(&room).Error; err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *chatRepository) JoinRoom(ctx context.Context, roomID, userID, role string) error {
	if err := r.db.WithContext(ctx).Save(&model.ChatMember{
		UserID: userID,
		RoomID: roomID,
		Role:   role,
	}).Error; err != nil {
		return err
	}
	return nil
}

func (r *chatRepository) JoinRoomWithTx(ctx context.Context, tx *gorm.DB, roomID, userID, role string) error {
	if err := tx.WithContext(ctx).Save(&model.ChatMember{
		UserID: userID,
		RoomID: roomID,
		Role:   role,
	}).Error; err != nil {
		return err
	}
	return nil
}

func (r *chatRepository) LeaveRoom(ctx context.Context, roomID, userID string) error {
	if err := r.db.WithContext(ctx).Where("room_id = ? AND user_id = ?", roomID, userID).
		Delete(&model.ChatMember{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *chatRepository) DeleteRoom(ctx context.Context, roomID string) error {
	if err := r.db.WithContext(ctx).Where("id = ?", roomID).
		Delete(&model.ChatRoom{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *chatRepository) ListVisibleRooms(ctx context.Context, userID string) ([]model.ChatRoom, error) {
	var rooms []model.ChatRoom
	if err := r.db.WithContext(ctx).Model(&model.ChatRoom{}).
		Joins("LEFT JOIN chat_room_member ON chat_room_member.room_id = chat_room.id AND chat_room_member.user_id = ?", userID).
		Where("chat_room.visibility = ? OR chat_room_member.user_id IS NOT NULL", VisibilityPublic).
		Find(&rooms).Error; err != nil {
		return nil, err
	}
	return rooms, nil
}

func (r *chatRepository) ListMyRooms(ctx context.Context, userID string) ([]model.ChatRoom, error) {
	var rooms []model.ChatRoom
	if err := r.db.WithContext(ctx).Model(&model.ChatRoom{}).
		Where("owner_id = ?", userID).Find(&rooms).Error; err != nil {
		return nil, err
	}
	return rooms, nil
}

func (r *chatRepository) FindMemberById(ctx context.Context, roomID, userID string) (*model.ChatMember, error) {
	var member model.ChatMember
	if err := r.db.WithContext(ctx).Where("room_id=? AND user_id = ?", roomID, userID).Take(&member).Error; err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *chatRepository) ListRoomMembers(ctx context.Context, roomID string) ([]model.ChatMember, error) {
	var members []model.ChatMember
	if err := r.db.WithContext(ctx).Model(&model.ChatMember{}).
		Where("room_id = ?", roomID).Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

func (r *chatRepository) ListRoomMessages(ctx context.Context, roomID string, page, pageSize int) ([]model.ChatMessage, error) {
	var messages []model.ChatMessage
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).Model(&model.ChatMessage{}).
		Where("room_id = ?", roomID).Order("created_at ASC").Limit(pageSize).Offset(offset).Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *chatRepository) CountRoomMessages(ctx context.Context, roomID string) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.ChatMessage{}).
		Where("room_id = ?", roomID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *chatRepository) MessagesHistory(ctx context.Context, roomID string, limit int) ([]model.ChatMessage, error) {
	var messages []model.ChatMessage
	if limit <= 0 || limit > 50 {
		limit = 50
	}
	if err := r.db.WithContext(ctx).Model(&model.ChatMessage{}).
		Where("room_id = ?", roomID).Order("created_at ASC").Limit(limit).Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *chatRepository) SaveMessage(ctx context.Context, message *model.ChatMessage) error {
	if err := r.db.WithContext(ctx).Create(message).Error; err != nil {
		return err
	}
	return nil
}
