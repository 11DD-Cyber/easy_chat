package conversation

import (
	"context"
	"easy_chat/apps/im/ws/internal/svc"
	"easy_chat/apps/im/ws/logic"
	"easy_chat/apps/im/ws/websocket"
	"easy_chat/apps/im/ws/ws"

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
		l := logic.NewConversation(context.Background(), srv, srvCtx)
		if err := l.SingleChat(&data, srv.GetUsers(conn)[0]); err != nil {
			srv.Send(websocket.NewErrMessage(err), conn.Conn)
			return
		}
		srv.SendByUserId(websocket.NewMessage(data.SendId, data), data.RecvId)

	}
}
