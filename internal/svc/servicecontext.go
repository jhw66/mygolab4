// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"log"

	"github.com/jhw66/myvideo_lab4/internal/config"
	"github.com/jhw66/myvideo_lab4/internal/middleware"
	"github.com/jhw66/myvideo_lab4/pkg/cache/cachemodel"
	"github.com/jhw66/myvideo_lab4/pkg/cache/cacherepo"
	"github.com/jhw66/myvideo_lab4/pkg/db/model"
	"github.com/jhw66/myvideo_lab4/pkg/db/repository/chat"
	"github.com/jhw66/myvideo_lab4/pkg/db/repository/comment"
	"github.com/jhw66/myvideo_lab4/pkg/db/repository/commentfavorite"
	"github.com/jhw66/myvideo_lab4/pkg/db/repository/favorite"
	"github.com/jhw66/myvideo_lab4/pkg/db/repository/relation"
	"github.com/jhw66/myvideo_lab4/pkg/db/repository/transaction"
	"github.com/jhw66/myvideo_lab4/pkg/db/repository/user"
	"github.com/jhw66/myvideo_lab4/pkg/db/repository/video"
	"github.com/jhw66/myvideo_lab4/pkg/webchat"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config                config.Config
	ChatHub               *webchat.Hub
	AccessAuth            rest.Middleware
	UserRepo              user.UserRepository
	VideoRepo             video.VideoRepository
	CommentRepo           comment.CommentRepository
	CommentFavoriteRepo   commentfavorite.CommentFavoriteRepository
	FavoriteRepo          favorite.FavoriteRepository
	RelationRepo          relation.RelationRepository
	TransactionRepository transaction.TransactionRepository
	ChatRepo              chat.ChatRepository
	CommentCache          cacherepo.CommentCache
	FavoriteCache         cacherepo.FavoriteCache
	RankCache             cacherepo.RankCache
	ChatCache             cacherepo.ChatCache
	//SessionAccount rest.Middleware
	//AuthLogin      rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	gormDB, err := model.InitDB(&c)
	if err != nil {
		log.Fatalf("init mysql failed:%v", err)
	}

	rdb, err := cachemodel.InitRedis(&c)
	if err != nil {
		log.Fatalf("init redis failed:%v", err)
	}

	userRepo := user.NewUserRepository(gormDB)

	return &ServiceContext{
		Config:                c,
		ChatHub:               webchat.NewHub(),
		AccessAuth:            middleware.NewAccessAuthMiddleware(userRepo, c).Handle,
		UserRepo:              userRepo,
		VideoRepo:             video.NewVideoRepository(gormDB),
		CommentRepo:           comment.NewCommentRepository(gormDB),
		CommentFavoriteRepo:   commentfavorite.NewCommentFavoriteRepository(gormDB),
		FavoriteRepo:          favorite.NewFavoriteRepository(gormDB),
		RelationRepo:          relation.NewRelationRepository(gormDB),
		TransactionRepository: transaction.NewTransactionRepository(gormDB),
		ChatRepo:              chat.NewChatRepository(gormDB),
		CommentCache:          cacherepo.NewCommentCache(rdb),
		FavoriteCache:         cacherepo.NewFavoriteCache(rdb),
		RankCache:             cacherepo.NewRankCache(rdb),
		ChatCache:             cacherepo.NewChatCache(rdb),
		//SessionAccount: middleware.NewSessionAccountMiddleware("").Handle,
		//AuthLogin:      middleware.NewAuthLoginMiddleware().Handle,
	}
}
