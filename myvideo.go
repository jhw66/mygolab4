// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/jhw66/myvideo_lab4/internal/config"
	"github.com/jhw66/myvideo_lab4/internal/handler"
	"github.com/jhw66/myvideo_lab4/internal/logic/core"
	"github.com/jhw66/myvideo_lab4/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/myvideo-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	if err := core.WarmUpRankZSet(context.Background(), ctx); err != nil {
		logx.Errorf("warmup rank zset failed: %v", err)
	}
	go core.StartVideoStatSync(context.Background(), ctx, core.DefaultSyncInterval, core.DefaultSyncBatchSize)
	go core.StartCommentFavoriteStatSync(context.Background(), ctx, core.DefaultSyncInterval, core.DefaultSyncBatchSize)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
