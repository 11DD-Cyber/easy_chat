// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package main

import (
	"flag"
	"fmt"

	"easy_chat/apps/social/api/internal/config"
	"easy_chat/apps/social/api/internal/handler"
	"easy_chat/apps/social/api/internal/svc"
	"easy_chat/pkg/resultx"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var configFile = flag.String("f", "etc/social.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)
	httpx.SetOkHandler(resultx.OkHandler)
	httpx.SetErrorHandlerCtx(resultx.ErrHandler(c.Name))
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
