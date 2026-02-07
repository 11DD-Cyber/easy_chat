package user

import (
	"easy_chat/apps/im/ws/internal/svc"
	"easy_chat/apps/im/ws/websocket"
)

func OnLine(svc *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		//获取所有在线用户
		uids := srv.GetUsers()
		//获取当前连接绑定用户
		u := srv.GetUsers(conn)
		if len(u) == 0 {
			srv.Logger.Errorf("上线处理失败：当前连接未绑定用户")
			return
		}

		err := srv.Send(websocket.NewMessage(u[0], uids), conn.Conn)
		if err != nil {
			srv.Logger.Errorf("error", err)
		}
	}
}
