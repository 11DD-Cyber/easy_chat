package svc

import (
	"easy_chat/apps/im/models"
	"easy_chat/apps/im/ws/internal/config"
	"easy_chat/apps/task/mq/mqclient"
)

type ServiceContext struct {
	Config config.Config
	models.ChatLogModel
	MsgChatTransfer mqclient.MsgChatTransferClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:          c,
		ChatLogModel:    models.NewChatLogModel(c.Mongo.Url, c.Mongo.Db, "chat_logs"),
		MsgChatTransfer: mqclient.NewMsgChatTransferClient(c.MsgChatTransfer.Addrs, c.MsgChatTransfer.Topic),
	}
}
