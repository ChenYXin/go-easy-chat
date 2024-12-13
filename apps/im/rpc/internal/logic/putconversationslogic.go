package logic

import (
	"context"
	"easy-chat/apps/im/immodels"
	"easy-chat/apps/im/rpc/im"
	"easy-chat/apps/im/rpc/internal/svc"
	"easy-chat/pkg/constants"
	"easy-chat/pkg/xerr"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type PutConversationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPutConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PutConversationsLogic {
	return &PutConversationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新会话
func (l *PutConversationsLogic) PutConversations(in *im.PutConversationsReq) (*im.PutConversationsResp, error) {
	// todo: add your logic here and delete this line

	// 根据用户查询用户会话列表（更新一定要有会话）
	data, err := l.svcCtx.ConversationsModel.FindByUserId(l.ctx, in.UserId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "Conversations FindByUserId err %v, req %v", err, in.UserId)
	}

	// 会话列表为空，创建一个空的
	if data.ConversationList == nil {
		data.ConversationList = make(map[string]*immodels.Conversation)
	}

	// 要更新的会话列表
	for s, conversion := range in.ConversationList {
		var oldTotal int
		// 总量不能为空，为空就不用更新了
		if data.ConversationList[s] != nil {
			oldTotal = data.ConversationList[s].Total
		}
		data.ConversationList[s] = &immodels.Conversation{
			ConversationId: conversion.ConversationId,
			ChatType:       constants.ChatType(conversion.ChatType),
			IsShow:         conversion.IsShow,
			Total:          int(conversion.Read) + oldTotal, // 已读总量，新读的+已经读过的
			Seq:            conversion.Seq,
		}

		// 更新数据库
		_, err = l.svcCtx.ConversationsModel.Update(l.ctx, data)
		if err != nil {
			return nil, errors.Wrapf(xerr.NewDBErr(), "Conversations Update err %v, req %v", err, data)
		}
	}

	return &im.PutConversationsResp{}, nil
}
