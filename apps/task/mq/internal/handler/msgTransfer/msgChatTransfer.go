package msgTransfer

import (
	"context"
	"easy-chat/apps/im/immodels"
	"easy-chat/apps/im/ws/ws"
	"easy-chat/apps/task/mq/internal/svc"
	"easy-chat/apps/task/mq/mq"
	"easy-chat/pkg/bitmap"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MsgChatTransfer struct {
	*baseMsgTransfer
}

func NewMsgChatTransfer(svc *svc.ServiceContext) *MsgChatTransfer {
	return &MsgChatTransfer{
		baseMsgTransfer: NewBaseMsgTransfer(svc),
	}
}

// Consume 实现消费消息的接口,go-zero提供的queue
func (m *MsgChatTransfer) Consume(ctx context.Context, key, value string) error {
	fmt.Println("key ：", key, " value : ", value)

	var (
		data  mq.MsgChatTransfer // 消息类型结构体
		c     = context.Background()
		msgId = primitive.NewObjectID()
	)

	// kafka保存的是json序列化数据，这里消费者拿到反序列化
	if err := json.Unmarshal([]byte(value), &data); err != nil {
		return err
	}

	// 消费者拿到数据之后把数据保存到数据库
	if err := m.addChatLog(c, msgId, &data); err != nil {
		return err
	}

	// 根据类型判断处理，从kafka中拿到数据消息，使用websocket发送
	return m.Transfer(ctx, &ws.Push{
		ConversationId: data.ConversationId,
		ChatType:       data.ChatType,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		RecvIds:        data.RecvIds,
		SendTime:       data.SendTime,
		Mtype:          data.Mtype,
		MsgId:          msgId.Hex(),
		Content:        data.Content,
	})
}

func (m *MsgChatTransfer) addChatLog(ctx context.Context, msgId primitive.ObjectID, data *mq.MsgChatTransfer) error {
	chatlog := immodels.ChatLog{
		ID:             msgId,
		ConversationId: data.ConversationId,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		ChatType:       data.ChatType,
		MsgFrom:        0,
		MsgType:        data.Mtype,
		MsgContent:     data.Content,
		SendTime:       data.SendTime,
	}

	// 记录已读未读情况
	readRecords := bitmap.NewBitmap(0)
	// 自己是已读的
	readRecords.Set(chatlog.SendId)
	chatlog.ReadRecords = readRecords.Export()

	// 插入数据库
	err := m.svcCtx.ChatLogModel.Insert(ctx, &chatlog)
	if err != nil {
		return err
	}

	// 更新会话 ，收到消息的同时，要更新会话
	return m.svcCtx.ConversationModel.UpdateMsg(ctx, &chatlog)
}
