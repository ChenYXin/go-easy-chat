package logic

import (
	"context"
	"easy-chat/apps/im/immodels"
	"easy-chat/pkg/xerr"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"

	"easy-chat/apps/im/rpc/im"
	"easy-chat/apps/im/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetConversationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConversationsLogic {
	return &GetConversationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取会话
func (l *GetConversationsLogic) GetConversations(in *im.GetConversationsReq) (*im.GetConversationsResp, error) {
	// todo: add your logic here and delete this line

	// 根据用户查询用户会话列表
	data, err := l.svcCtx.ConversationsModel.FindByUserId(l.ctx, in.UserId)
	if err != nil {
		if errors.Is(err, immodels.ErrNotFound) { // 可能会话是空的，也就是没有聊天记录
			return &im.GetConversationsResp{}, nil
		}
		return nil, errors.Wrapf(xerr.NewDBErr(), "Conversations FindByUserId err %v, req %v", err, in.UserId)
	}
	var res im.GetConversationsResp
	copier.Copy(&res, &data)

	// 根据会话列表，查询具体的会话
	ids := make([]string, 0, len(data.ConversationList))
	for _, conversion := range data.ConversationList {
		ids = append(ids, conversion.ConversationId)
	}
	conversations, err := l.svcCtx.ConversationModel.ListByConversationIds(l.ctx, ids)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "ConversationModel ListByConversationIds err %v, req %v", err, ids)
	}

	// 计算是否存在未读消息
	for _, conversation := range conversations {
		// 判断会话列表中有没有存在这样的会话
		if _, ok := res.ConversationList[conversation.ConversationId]; !ok {
			continue
		}

		// 用户读取的消息量
		total := res.ConversationList[conversation.ConversationId].Total
		// 用户读取的消息量比会话的总消息量少，说明有未读的消息
		if total < int32(conversation.Total) {
			// 消息总量更新，不更新数据库
			res.ConversationList[conversation.ConversationId].Total = int32(conversation.Total)
			// 有多少消息是未读的(会话总 - 当前读的)
			res.ConversationList[conversation.ConversationId].ToRead = int32(conversation.Total) - total
			// 更改当前会话为显示状态(可能没有会话，在聊天列表中就被删除对话框的)
			res.ConversationList[conversation.ConversationId].IsShow = true
		}
	}

	return &res, nil
}
