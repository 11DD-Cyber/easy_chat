package handler

import (
	"easy_chat/apps/im/ws/internal/handler/conversation"
	"easy_chat/apps/im/ws/internal/handler/user"
	"easy_chat/apps/im/ws/internal/svc"
	"easy_chat/apps/im/ws/websocket"
)

func RegisterHandlers(srv *websocket.Server, svc *svc.ServiceContext) {
	srv.AddRoutes([]websocket.Route{
		{
			Method:  "user.onLine",
			Handler: user.OnLine(svc),
		},
		{
			Method:  "conversation.chat",
			Handler: conversation.Chat(svc),
		},
	})
}
