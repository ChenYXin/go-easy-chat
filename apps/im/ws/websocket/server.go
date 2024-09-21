package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"sync"
)

type Server struct {
	sync.RWMutex

	opt            *serverOptions
	authentication Authentication

	routes     map[string]HandlerFunc
	addr       string
	patten     string
	connToUser map[*Conn]string
	userToConn map[string]*Conn

	upgrader websocket.Upgrader
	logx.Logger
}

func NewServer(addr string, opts ...ServerOptions) *Server {
	opt := newServerOptions(opts...)
	return &Server{
		routes:   make(map[string]HandlerFunc),
		addr:     addr,
		patten:   opt.patten,
		opt:      &opt,
		upgrader: websocket.Upgrader{},

		authentication: opt.Authentication,

		connToUser: make(map[*Conn]string),
		userToConn: make(map[string]*Conn),

		Logger: logx.WithContext(context.Background()),
	}
}

func (s *Server) ServerWs(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			s.Errorf("server handler ws recover err : %v", r)
		}
	}()

	//获取连接对象
	conn := NewConn(s, w, r)
	if conn == nil {
		return
	}
	//conn, err := s.upgrader.Upgrade(w, r, nil)
	//if err != nil {
	//	fmt.Println("upgrade:", err)
	//	s.Errorf("server upgrade err : %v", err)
	//	return
	//}
	if !s.authentication.Auth(w, r) {
		//conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("不具备访问权限")))
		s.Send(&Message{FrameType: FrameData, Data: fmt.Sprintf("不具备访问权限")}, conn)
		conn.Close()
		return
	}
	//记录连接
	s.addConn(conn, r)
	//处理连接
	go s.handlerConn(conn)
}

func (s *Server) addConn(conn *Conn, req *http.Request) {
	uid := s.authentication.UserId(req)

	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	//验证用户是否之前登陆过
	if c := s.userToConn[uid]; c != nil {
		//关闭之前的连接
		fmt.Println("验证用户是否之前登陆过")
		c.Close()
	}
	s.connToUser[conn] = uid
	s.userToConn[uid] = conn
}

func (s *Server) GetConn(uid string) *Conn {
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()
	return s.userToConn[uid]
}

func (s *Server) GetConns(uids ...string) []*Conn {
	if len(uids) == 0 {
		return nil
	}
	s.RWMutex.RLock()
	defer s.RUnlock()
	res := make([]*Conn, 0, len(uids))
	for _, uid := range uids {
		res = append(res, s.userToConn[uid])
	}
	return res
}

func (s *Server) GetUsers(conns ...*Conn) []string {
	s.RWMutex.RLock()
	defer s.RUnlock()
	//return s.userToConn[uid]
	var res []string
	if len(conns) == 0 {
		//获取全部
		res = make([]string, 0, len(s.connToUser))
		for _, uid := range s.connToUser {
			res = append(res, uid)
		}
	} else {
		//获取部分
		res = make([]string, 0, len(conns))
		for _, conn := range conns {
			res = append(res, s.connToUser[conn])
		}
	}
	return res
}

func (s *Server) Close(conn *Conn) {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	uid := s.connToUser[conn]
	if uid == "" {
		//已经被关闭了
		return
	}

	delete(s.connToUser, conn)
	delete(s.userToConn, uid)
	conn.Close()
}

func (s *Server) SendByUserId(msg interface{}, sendIds ...string) error {
	if len(sendIds) == 0 {
		return nil
	}
	return s.Send(msg, s.GetConns(sendIds...)...)
}
func (s *Server) Send(msg interface{}, conns ...*Conn) error {
	if len(conns) == 0 {
		return nil
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	for _, conn := range conns {
		if err = conn.WriteMessage(websocket.TextMessage, data); err != nil {
			return err
		}
	}
	return nil
}

// 根据连接对象,执行任务处理
func (s *Server) handlerConn(conn *Conn) {
	uids := s.GetUsers(conn)
	conn.Uid = uids[0]
	for {
		// 获取请求消息
		_, msg, err := conn.ReadMessage()
		if err != nil {
			s.Errorf("server read msg err : %v", err)
			//t关闭连接对象
			s.Close(conn)
			return
		}
		//解析消息
		var message Message
		if err = json.Unmarshal(msg, &message); err != nil {
			s.Errorf("json Unmarshal err : %v , msg : %v", err, msg)
			// 关闭连接对象
			s.Close(conn)
			return
		}

		//根据消息进行处理
		switch message.FrameType {
		case FramePing:
			s.Send(&Message{FrameType: FramePing}, conn)
		case FrameData:
			//根据请求的method 分发路由并执行
			if handler, ok := s.routes[message.Method]; ok {
				handler(s, conn, &message)
			} else {
				//conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("不存在执行的方法 %v 请检查", message.Method)))
				s.Send(&Message{
					FrameType: FrameData,
					Data:      fmt.Sprintf("不存在执行的方法 %v 请检查", message.Method),
				}, conn)

			}
		}
	}

}

func (s *Server) AddRouters(rs []Route) {
	for _, r := range rs {
		s.routes[r.Method] = r.Handler
	}
}

func (s *Server) Start() {
	http.HandleFunc(s.patten, s.ServerWs)
	fmt.Println("websocket server start")
	s.Info(http.ListenAndServe(s.addr, nil))
}

func (s *Server) Stop() {
	fmt.Println("websocket server stop")

}
