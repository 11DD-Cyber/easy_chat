package msgtransfer

import (
	"context"
	"easy_chat/apps/im/models"
	"easy_chat/apps/im/ws/websocket"
	"easy_chat/apps/im/ws/ws"
	"easy_chat/apps/task/mq/internal/svc"
	"easy_chat/apps/task/mq/mq"
	"easy_chat/pkg/constants"
	"encoding/json"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/logx"
)

type MsgChatTransfer struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewMsgChatTransfer(svc *svc.ServiceContext) kq.ConsumeHandler {
	return &MsgChatTransfer{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svc,
	}
}

func (m *MsgChatTransfer) Consume(ctx context.Context, key, value string) error {
	var data mq.MsgChatTransfer
	if err := json.Unmarshal([]byte(value), &data); err != nil {
		return err
	}
	//记录消息
	if err := m.addChatLog(ctx, data); err != nil {
		return err
	}
	//推送发送
	return m.svcCtx.WsClient.Send(websocket.Message{
		FrameType: websocket.FrameData,
		Method:    "push",
		FromId:    constants.SYSTEM_ROOT_UID,
		Data: ws.Push{
			ConversationId: data.ConversationId,
			SendId:         data.SendId,
			RecvId:         data.RecvId,
			SendTime:       data.SendTime,
			Msg: ws.Msg{
				MType:   data.Msgtype,
				Content: data.Content,
			},
		},
	})

}

func (m *MsgChatTransfer) addChatLog(ctx context.Context, data mq.MsgChatTransfer) error {
	chatLog := models.ChatLog{
		ConversationId: data.ConversationId,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		MsgType:        data.Msgtype,
		MsgContent:     data.Content,
		SendTime:       data.SendTime,
	}
	err := m.svcCtx.ChatLogModel.Insert(ctx, &chatLog)
	if err != nil {
		return err
	}
	return m.svcCtx.ConversationModel.UpdateMsg(ctx, &chatLog)
}
