package serializer

import (
	"github.com/jhw66/myvideo_lab4/internal/types"
	"github.com/jhw66/myvideo_lab4/pkg/db/model"
)

func ChatRoomItemFromModel(room *model.ChatRoom) *types.ChatRoomItem {
	return &types.ChatRoomItem{
		ID:         room.ID,
		RoomName:   room.RoomName,
		Visibility: room.Visibility,
		OwnerID:    room.OwnerID,
		CreatedAt:  room.CreatedAt.Unix(),
	}
}

func ChatRoomListRspFromModels(rooms []model.ChatRoom) *types.GetRoomResp {
	items := make([]types.ChatRoomItem, 0, len(rooms))
	for i := range rooms {
		items = append(items, *ChatRoomItemFromModel(&rooms[i]))
	}
	return &types.GetRoomResp{
		Status: 200,
		Data:   &items,
		Msg:    "获取房间列表成功",
		Error:  "",
	}
}

func ChatMessageItemFromModel(message *model.ChatMessage) *types.ChatMessageItem {
	return &types.ChatMessageItem{
		ID:        message.ID,
		RoomID:    message.RoomID,
		UserID:    message.UserID,
		Content:   message.Content,
		CreatedAt: message.CreatedAt.Unix(),
	}
}

func ChatMessageListRspFromModels(messages []model.ChatMessage, total int64, page int, pageSize int) *types.GetMessagesResp {
	items := make([]types.ChatMessageItem, 0, len(messages))
	for i := range messages {
		items = append(items, *ChatMessageItemFromModel(&messages[i]))
	}
	return &types.GetMessagesResp{
		Status: 200,
		Data: &types.ChatMessageListData{
			Total:    total,
			Page:     page,
			PageSize: pageSize,
			Messages: &items,
		},
		Msg:   "获取消息成功",
		Error: "",
	}
}
