package msgTransfer

import (
	"context"
	"easy-chat/apps/im/immodels"
	"easy-chat/apps/im/ws/websocket"
	"easy-chat/apps/task/mq/internal/svc"
	"easy-chat/apps/task/mq/mq"
	"easy-chat/pkg/constants"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
)

type MsgChatTransfer struct {
	logx.Logger
	svc *svc.ServiceContext
}

func NewMsgChatTransfer(svc *svc.ServiceContext) *MsgChatTransfer {
	return &MsgChatTransfer{
		Logger: logx.WithContext(context.Background()),
		svc:    svc,
	}
}

// Consume 实现消费消息的接口,go-zero提供的queue
func (m *MsgChatTransfer) Consume(ctx context.Context, key, value string) error {
	fmt.Println("key：", key, " value：", value)

	var (
		data mq.MsgChatTransfer
		//ctx  = context.Background()
	)

	if err := json.Unmarshal([]byte(value), &data); err != nil {
		return err
	}

	//记录消息，存储到mangodb
	if err := m.addChatLog(ctx, &data); err != nil {
		return err
	}

	//推送消息
	return m.svc.WsClient.Send(websocket.Message{
		FrameType: websocket.FrameData,
		Method:    "push",
		FormId:    constants.SYSTEM_ROOT_UID,
		Data:      data,
	})
}

func (m *MsgChatTransfer) addChatLog(ctx context.Context, data *mq.MsgChatTransfer) error {
	//记录消息，存储到mangodb
	chatLog := immodels.ChatLog{
		ConversationId: data.ConversationId,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		MsgFrom:        0,
		ChatType:       data.ChatType,
		MsgType:        data.Mtype,
		MsgContent:     data.Content,
		SendTime:       data.SendTime,
	}
	return m.svc.ChatLogModel.Insert(ctx, &chatLog)
}
