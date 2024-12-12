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
		fmt.Println("消息内容：", data)

		fmt.Println("发送消息的类型1:群聊天，2:单聊 --- ", data.ChatType)
		switch data.ChatType {
		case constants.SingleChatType:
			//将消息推送给kafka
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
				fmt.Println("ws发送消息给kafka发生错误", err.Error())
				err = srv.Send(websocket.NewErrMessage(err), conn)
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
