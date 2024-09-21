package msgTransfer

import (
	"context"
	"easy-chat/apps/task/mq/internal/svc"
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

// 实现消费消息的接口
func (m *MsgChatTransfer) Consume(ctx context.Context, key, value string) error {
	fmt.Println("key：", key, " value：", value)
	return nil
}
