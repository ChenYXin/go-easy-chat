package conversation

import (
	"easy-chat/apps/im/ws/internal/svc"
	"easy-chat/apps/im/ws/websocket"
	"easy-chat/apps/im/ws/ws"
	"easy-chat/apps/task/mq/mq"
	"easy-chat/pkg/constants"
	"fmt"
	"github.com/mitchellh/mapstructure"
)

func Chat(svc *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		//todo : 私聊
		var data ws.Chat
		if err := mapstructure.Decode(msg.Data, &data); err != nil {
			fmt.Println(err.Error())
			srv.Send(websocket.NewErrMessage(err), conn)
			return
		}

		switch data.ChatType {
		case constants.SingleChatType:

			err := svc.MsgChatTransferClient.Push(&mq.MsgChatTransfer{
				ConversationId: data.ConversationId,
				ChatType:       data.ChatType,
				SendId:         data.SendId,
				RecvId:         data.RecvId,
				SendTime:       data.SendTime,
				Mtype:          data.Msg.Mtype,
				Content:        data.Msg.Content,
			})
			if err != nil {
				srv.Send(websocket.NewErrMessage(err), conn)
				return
			}

			//err := logic.NewConversation(context.Background(), srv, svc).SingleChat(&data, conn.Uid)
			//if err != nil {
			//	fmt.Println(err.Error())
			//	srv.Send(websocket.NewErrMessage(err), conn)
			//	return
			//}
			//srv.SendByUserId(websocket.NewMessage(conn.Uid, ws.Chat{
			//	ConversationId: data.ConversationId,
			//	ChatType:       data.ChatType,
			//	SendId:         conn.Uid,
			//	RecvId:         data.RecvId,
			//	SendTime:       time.Now().UnixMilli(),
			//	Msg:            data.Msg,
			//}), data.RecvId)
		default:
		}
	}
}
