package msgTransfer

import (
	"context"
	"easy-chat/apps/im/ws/websocket"
	"easy-chat/apps/im/ws/ws"
	"easy-chat/apps/social/rpc/socialclient"
	"easy-chat/apps/task/mq/internal/svc"
	"easy-chat/pkg/constants"
	"github.com/zeromicro/go-zero/core/logx"
)

type baseMsgTransfer struct {
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBaseMsgTransfer(svc *svc.ServiceContext) *baseMsgTransfer {
	return &baseMsgTransfer{
		svcCtx: svc,
		Logger: logx.WithContext(context.Background()),
	}
}

// Transfer
//
//	@Description: mq转发消息的
//	@receiver m
//	@return error
func (m *baseMsgTransfer) Transfer(ctx context.Context, data *ws.Push) error {
	// 私聊和群聊区分
	var err error
	switch data.ChatType {
	case constants.GroupChatType:
		err = m.group(ctx, data)
	case constants.SingleChatType:
		err = m.single(ctx, data)
	}
	return err
}

func (m *baseMsgTransfer) single(ctx context.Context, data *ws.Push) error {
	return m.svcCtx.WsClient.Send(websocket.Message{
		FrameType: websocket.FrameData,
		Method:    "push",
		FormId:    constants.SYSTEM_ROOT_UID, // 当前系统角色
		Data:      data,
	})
}

func (m *baseMsgTransfer) group(ctx context.Context, data *ws.Push) error {
	// 查询群用户，推送给对应的群用户, 不是从数据库中，直接用social rpc服务获取指定的群用户
	users, err := m.svcCtx.Social.GroupUsers(ctx, &socialclient.GroupUsersReq{
		GroupId: data.RecvId, // recvId就是群Id
	})
	if err != nil {
		return err
	}

	data.RecvIds = make([]string, 0, len(users.List))

	for _, user := range users.List {
		// 不能发送给自己
		if user.UserId == data.SendId {
			continue
		}

		data.RecvIds = append(data.RecvIds, user.UserId)
	}

	// 调到了多协程发送
	return m.svcCtx.WsClient.Send(websocket.Message{
		FrameType: websocket.FrameData,
		Method:    "push",
		FormId:    constants.SYSTEM_ROOT_UID, // 当前系统角色
		Data:      data,
	})
}
