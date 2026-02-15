package main

import (
	"easy_chat/apps/task/mq/internal/config"
	"easy_chat/apps/task/mq/internal/handler"
	"easy_chat/apps/task/mq/internal/svc"
	"flag"
	"fmt"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
)

var configFile = flag.String("f", "etc/task.yaml", "the config file")

func main() {
	flag.Parse()
	// 加载配置文件
	var c config.Config
	conf.MustLoad(*configFile, &c)
	if err := c.SetUp(); err != nil {
		panic(err)
	}
	serviceGroup := service.NewServiceGroup()
	defer serviceGroup.Stop()
	svcCtx := svc.NewServiceContext(c)
	listen := handler.NewListen(svcCtx)
	for _, s := range listen.Services() {
		serviceGroup.Add(s)
	}
	fmt.Println("Starting mq server...")
	serviceGroup.Start()

}
