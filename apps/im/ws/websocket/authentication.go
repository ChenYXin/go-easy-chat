package websocket

import (
	"fmt"
	"net/http"
	"time"
)

type Authentication interface {
	Auth(w http.ResponseWriter, r *http.Request) bool //鉴权是否成功
	UserId(r *http.Request) string                    //获取用户的ID
}

type authentication struct {
}

func (a *authentication) Auth(w http.ResponseWriter, r *http.Request) bool {
	return true
}

func (a *authentication) UserId(r *http.Request) string {
	query := r.URL.Query()
	//如果请求的参数携带userId,则直接返回
	if query != nil && query["userId"] != nil {
		return fmt.Sprintf("%v", query["userId"])
	}
	//返则使用时间戳作为userId
	return fmt.Sprintf("%v", time.Now().UnixMilli())
}
