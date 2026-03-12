package push

import (
	"context"

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
		chat := &ws.Chat{
			ConversationId: data.ConversationId,
			SendId:         data.SendId,
			RecvId:         data.RecvId,
			SendTime:       data.SendTime,
			Msg: ws.Msg{
				MType:   data.Msg.MType,
				Content: data.Msg.Content,
			},
		}
		rconn := srv.GetConns(data.RecvId)
		if len(rconn) == 0 {
			if err := svcCtx.SaveOfflineMessage(context.Background(), data.RecvId, chat); err != nil {
				srv.Logger.Errorf("save offline message err recvId %s: %v", data.RecvId, err)
			}
			return
		}
		if err := srv.Send(websocket.NewMessage(chat.SendId, chat), rconn...); err != nil {
			srv.Logger.Errorf("send live message err recvId %s: %v", data.RecvId, err)
		}
	}
}
