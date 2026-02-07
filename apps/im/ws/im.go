package main

import (
	"easy_chat/apps/im/ws/internal/config"
	"easy_chat/apps/im/ws/internal/handler"
	"easy_chat/apps/im/ws/internal/svc"
	server "easy_chat/apps/im/ws/websocket"
	"flag"
	"fmt"

	"github.com/zeromicro/go-zero/core/conf"
)

var configFile = flag.String("f", "etc/dev/im.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	if err := c.SetUp(); err != nil {
		panic(err)
	}
	ctx := svc.NewServiceContext(c)
	srv := server.NewServer(c.ListenOn, handler.NewJwtAuth(ctx))
	defer srv.Stop()

	handler.RegisterHandlers(srv, ctx)

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	srv.Start()
}
