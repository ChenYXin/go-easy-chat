package handler

import (
	"context"
	"easy-chat/apps/im/ws/internal/svc"
	"easy-chat/pkg/ctxdata"
	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/token"
	"net/http"
)

// JwtAuth 在im中需要实现自定义的authebtication对象
type JwtAuth struct {
	svc    *svc.ServiceContext
	parser *token.TokenParser //解析token的对象
	logx.Logger
}

func NewJwtAuth(svc *svc.ServiceContext) *JwtAuth {
	return &JwtAuth{
		svc:    svc,
		parser: token.NewTokenParser(),
		Logger: logx.WithContext(context.Background()),
	}
}

// Auth 鉴权的过程可以参考go-zero的authhandler
func (j *JwtAuth) Auth(w http.ResponseWriter, r *http.Request) bool {
	//解析token
	tok, err := j.parser.ParseToken(r, j.svc.Config.JwtAuth.AccessSecret, "")
	if err != nil {
		j.Errorf("parse token err: %v", err)
		return false
	}
	//验证token是否合法
	if !tok.Valid {
		return false
	}
	//获取token里面的信息
	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}
	//将token里面的数据放入到request中
	*r = *r.WithContext(context.WithValue(r.Context(), ctxdata.Identify, claims[ctxdata.Identify]))
	return true
}

func (j *JwtAuth) UserId(r *http.Request) string {
	return ctxdata.GetUId(r.Context())
}
