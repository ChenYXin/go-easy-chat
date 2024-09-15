package user

import (
	"context"
	"easy-chat/apps/user/api/internal/svc"
	"easy-chat/apps/user/api/internal/types"
	"easy-chat/apps/user/rpc/user"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户登陆
func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	loginResp, err := l.svcCtx.User.Login(l.ctx, &user.LoginRequest{
		Phone:    req.Phone,
		Password: req.Password,
	})

	if err != nil {
		return nil, err
	}
	var res types.LoginResp
	copier.Copy(&res, loginResp)

	return &res, nil
}
