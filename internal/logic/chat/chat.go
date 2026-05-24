// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package chat

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jhw66/myvideo_lab4/internal/svc"
	"github.com/jhw66/myvideo_lab4/internal/types"
	"github.com/jhw66/myvideo_lab4/pkg/auth"
	"github.com/jhw66/myvideo_lab4/pkg/webchat"
	"github.com/zeromicro/go-zero/core/logx"
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type ChatLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatLogic {
	return &ChatLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChatLogic) Chat(w http.ResponseWriter, r *http.Request, req *types.ChatReq) (err error) {
	user, ok := auth.GetUserFromContext(l.ctx)
	if !ok || user == nil {
		return errors.New("用户未登录")
	}
	if member, err := l.svcCtx.ChatRepo.FindMemberById(l.ctx, req.RoomID, user.ID); member == nil || err != nil {
		return errors.New("您没有权限访问该房间或者访问失败")
	}

	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade to websocket failed:", err)
		return err
	}

	client := webchat.NewClient(l.svcCtx.ChatHub, conn, req.RoomID, user.ID, user.UserName)
	client.ServeClient(l.ctx, conn, l.svcCtx.ChatRepo)

	return nil
}
