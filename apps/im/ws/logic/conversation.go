package logic

import (
	"context"
	"easy_chat/apps/im/models"
	"easy_chat/apps/im/ws/internal/svc"
	"easy_chat/apps/im/ws/websocket"
	"easy_chat/apps/im/ws/ws"
	"easy_chat/pkg/wuid"
	"time"
)

type Conversation struct {
	ctx context.Context
	srv *websocket.Server
	svc *svc.ServiceContext
}

func NewConversation(ctx context.Context, srv *websocket.Server, svc *svc.ServiceContext) *Conversation {
	return &Conversation{
		ctx: ctx,
		srv: srv,
		svc: svc,
	}
}
func (l *Conversation) SingleChat(data *ws.Chat, userId string) error {
	if data.ConversationId == "" {
		data.ConversationId = wuid.CombineId(userId, data.RecvId)
	}
	//记录消息
	chatLog := models.ChatLog{
		ConversationId: data.ConversationId,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		MsgFrom:        0,
		MsgType:        data.Msg.MType,
		MsgContent:     data.Msg.Content,
		SendTime:       time.Now().UnixMilli(),
		Status:         0,
	}
	err := l.svc.ChatLogModel.Insert(l.ctx, &chatLog)
	return err
}
