package push

import (
	"easy_chat/apps/im/ws/internal/svc"
	"easy_chat/apps/im/ws/websocket"
	"easy_chat/apps/im/ws/ws"

	"github.com/go-viper/mapstructure/v2"
)

func Push(svcCtx *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		var data ws.Push
		if err := mapstructure.Decode(msg.Data, &data); err != nil {
			srv.Send(websocket.NewErrMessage(err), conn.Conn)
			return
		}
		rconn := srv.GetConns(data.RecvId)
		if rconn == nil {
			//离线
			return
		}
		srv.Logger.Infof("解析消息失败：原始数据：%v", msg.Data)
		srv.Send(websocket.NewMessage(data.SendId, &ws.Chat{
			ConversationId: data.ConversationId,
			Msg: ws.Msg{
				MType:   data.Msg.MType,
				Content: data.Msg.Content,
			},
		}), rconn...)
	}
}
