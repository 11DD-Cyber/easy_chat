package conversation

import (
	"easy_chat/apps/im/ws/internal/svc"
	"easy_chat/apps/im/ws/websocket"
	"easy_chat/apps/im/ws/ws"
	"easy_chat/apps/task/mq/mq"
	"easy_chat/pkg/wuid"
	"time"

	"github.com/go-viper/mapstructure/v2"
)

func Chat(srvCtx *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		//解析消息
		var data ws.Chat
		if err := mapstructure.Decode(msg.Data, &data); err != nil {
			srv.Send(websocket.NewErrMessage(err), conn.Conn)
			return
		}
		if data.ConversationId == "" {
			data.ConversationId = wuid.CombineId(data.SendId, data.RecvId)
		}
		err := srvCtx.MsgChatTransfer.Push(&mq.MsgChatTransfer{
			ConversationId: data.ConversationId,
			SendId:         data.SendId,
			RecvId:         data.RecvId,
			Msgtype:        data.Msg.MType,
			Content:        data.Msg.Content,
			SendTime:       time.Now().UnixNano(),
		})
		if err != nil {
			srv.Send(websocket.NewErrMessage(err), conn.Conn)
		}
	}
}
