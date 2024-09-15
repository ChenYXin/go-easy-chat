package user

import (
	"context"
	"easy-chat/apps/user/api/internal/svc"
	"easy-chat/apps/user/api/internal/types"
	"easy-chat/apps/user/rpc/user"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户注册
func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {

	registerRsp, err := l.svcCtx.User.Register(l.ctx, &user.RegisterRequest{
		Phone:    req.Phone,
		Nickname: req.Nickname,
		Password: req.Password,
		Avatar:   req.Avatar,
		Sex:      int32(req.Sex),
	})

	if err != nil {
		return nil, err
	}
	var res types.RegisterResp
	copier.Copy(&res, registerRsp)

	return &res, nil
}
