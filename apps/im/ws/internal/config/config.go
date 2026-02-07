package config

import "github.com/zeromicro/go-zero/core/service"

type Config struct {
	// 嵌入ServiceConf：复用go-zero的通用服务配置
	// 包含服务名（Name）、日志、监控、链路追踪、运行模式（dev/prod）等通用配置项
	service.ServiceConf
	ListenOn string
	JwtAuth  struct {
		AccessSecret string
	}
	Mongo struct {
		Url string
		Db  string
	}
}
