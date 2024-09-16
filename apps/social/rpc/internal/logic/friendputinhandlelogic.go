package logic

import (
	"context"
	"easy-chat/apps/social/socialmodels"
	"easy-chat/pkg/constants"
	"easy-chat/pkg/xerr"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"easy-chat/apps/social/rpc/internal/svc"
	"easy-chat/apps/social/rpc/social"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrFriendReqBeforePass   = xerr.NewMsg("好友申请并已经通过")
	ErrFriendReqBeforeRefuse = xerr.NewMsg("好友申请已经被拒绝")
)

type FriendPutInHandleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendPutInHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInHandleLogic {
	return &FriendPutInHandleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FriendPutInHandleLogic) FriendPutInHandle(in *social.FriendPutInHandleReq) (*social.FriendPutInHandleResp, error) {
	//获取好友申请记录
	friendReq, err := l.svcCtx.FriendRequestsModel.FindOne(l.ctx, int64(in.FriendReqId))
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find friendRequest err %v req %v", err, in.FriendReqId)
	}
	//验证是否有处理
	switch constants.HandlerResult(friendReq.HandleResult.Int64) {
	case constants.PassHandlerResult:
		return nil, errors.WithStack(ErrFriendReqBeforePass)
	case constants.RefuseHandlerResult:
		return nil, errors.WithStack(ErrFriendReqBeforeRefuse)
	}
	fmt.Println("-----1")
	friendReq.HandleResult.Int64 = int64(in.HandleResult)
	fmt.Println("-----2")
	//修改申请结果 -》 通过【建立2条好友关系记录】 -》 事务
	err = l.svcCtx.FriendRequestsModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		fmt.Println("-----3")
		if err := l.svcCtx.FriendRequestsModel.Update(l.ctx, session, friendReq); err != nil {

			return errors.Wrapf(xerr.NewDBErr(), "update friend request err %v req %v", err, in.FriendReqId)
		}
		fmt.Println("-----4")
		if constants.HandlerResult(in.HandleResult) != constants.PassHandlerResult {
			return nil
		}
		fmt.Println("-----5")
		friends := []*socialmodels.Friends{
			{UserId: friendReq.UserId, FriendUid: friendReq.ReqUid},
			{UserId: friendReq.ReqUid, FriendUid: friendReq.UserId},
		}
		fmt.Println("-----6")
		_, err = l.svcCtx.FriendsModel.Inserts(l.ctx, session, friends...)
		if err != nil {
			fmt.Println("-----", err)
			return errors.Wrapf(xerr.NewDBErr(), "friend insert  request err %v", err)
		}
		fmt.Println("-----7")
		return nil
	})
	return &social.FriendPutInHandleResp{}, nil
}
