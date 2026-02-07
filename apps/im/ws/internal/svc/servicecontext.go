package svc

import (
	"easy_chat/apps/im/models"
	"easy_chat/apps/im/ws/internal/config"
)

type ServiceContext struct {
	Config config.Config
	models.ChatLogModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:       c,
		ChatLogModel: models.NewChatLogModel(c.Mongo.Url, c.Mongo.Db, "chat_logs"),
	}
}
