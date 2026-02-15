package svc

import (
	"easy_chat/apps/im/models"
	"easy_chat/apps/im/rpc/internal/config"
)

type ServiceContext struct {
	Config config.Config
	models.ChatLogModel
	models.ConversationModel
	models.ConversationsModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:             c,
		ChatLogModel:       models.NewChatLogModel(c.Mongo.Url, c.Mongo.Db, "chat_logs"),
		ConversationModel:  models.NewConversationModel(c.Mongo.Url, c.Mongo.Db, "chat_logs"),
		ConversationsModel: models.NewConversationsModel(c.Mongo.Url, c.Mongo.Db, "chat_logs"),
	}
}
