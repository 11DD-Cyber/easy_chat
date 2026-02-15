package svc

import (
	"easy_chat/apps/im/models"
	"easy_chat/apps/im/ws/websocket"
	"easy_chat/apps/task/mq/internal/config"
	"easy_chat/pkg/constants"
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

type ServiceContext struct {
	Config   config.Config
	WsClient websocket.Client
	Redis    *redis.Redis
	models.ChatLogModel
	models.ConversationModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	svcCtx := &ServiceContext{
		Config:            c,
		Redis:             redis.MustNewRedis(c.Redisx),
		ChatLogModel:      models.NewChatLogModel(c.Mongo.Url, c.Mongo.Db, "chat_logs"),
		ConversationModel: models.NewConversationModel(c.Mongo.Url, c.Mongo.Db, "chat_logs"),
	}
	token, err := svcCtx.GetToken()
	if err != nil {
		fmt.Printf("getToken err %v", err)
	}
	header := http.Header{}
	header.Set("Authorization", token)
	svcCtx.WsClient = websocket.NewClient(c.Ws.Host, "/ws", header)
	return svcCtx
}

func (svcCtx *ServiceContext) GetToken() (string, error) {
	return svcCtx.Redis.Get(constants.REDIS_SYSTEM_ROOT_TOKEN)
}
