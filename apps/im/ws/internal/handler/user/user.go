package user

import (
	"easy-chat/apps/im/ws/internal/svc"
	"easy-chat/apps/im/ws/websocket"
)

func OnLine(svc *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		uids := srv.GetUsers()

		u := srv.GetUsers(conn)
		//根据连接对象，获取在线的人数
		err := srv.Send(websocket.NewMessage(u[0], uids), conn)
		srv.Info("err", err)
	}
}
