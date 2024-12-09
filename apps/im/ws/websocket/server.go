package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"sync"
	"time"
)

type AckType int

const (
	NoAck AckType = iota
	OnlyAck
	RigorAck
)

func (t AckType) ToString() string {
	switch t {
	case OnlyAck:
		return "OnlyAck"
	case RigorAck:
		return "rigorAck"
	}
	return "NoAck"
}

type Server struct {
	sync.RWMutex // connToUser、userToConn 会有高并发的问题，需要加锁，保证线程安全

	opt            *serverOptions
	authentication Authentication //鉴权

	routes     map[string]HandlerFunc //存储实际具体要执行的路由，key就是方法名，value就是具体的执行方法
	addr       string
	patten     string
	connToUser map[*Conn]string //根据连接对象获取用户
	userToConn map[string]*Conn //根据用户获取连接对象

	upgrader websocket.Upgrader
	logx.Logger
}

func NewServer(addr string, opts ...ServerOptions) *Server {
	opt := newServerOptions(opts...)
	return &Server{
		routes:   make(map[string]HandlerFunc),
		addr:     addr,
		patten:   opt.pattern,
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
	//连接鉴权
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

// 添加连接对象
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
	//存在map中，key:连接对象,value:用户ID
	s.connToUser[conn] = uid
	//存在map中，key:用户ID,value:连接对象
	s.userToConn[uid] = conn
}

// GetConn 根据uid获得连接对象
func (s *Server) GetConn(uid string) *Conn {
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()
	return s.userToConn[uid]
}

// GetConns 根据多个uid获得多个连接对象
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

// GetUsers 根据连接对象获得userId
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

// Close 关闭连接对象
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

// SendByUserId 根据sendIds发送消息
func (s *Server) SendByUserId(msg interface{}, sendIds ...string) error {
	if len(sendIds) == 0 {
		return nil
	}
	return s.Send(msg, s.GetConns(sendIds...)...)
}

// Send 根据连接对象发送消息
func (s *Server) Send(msg interface{}, conns ...*Conn) error {
	if len(conns) == 0 {
		return nil
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	for _, conn := range conns {
		//通过websocket发送消息
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

	//处理任务
	go s.handlerWrite(conn)

	if s.isAck(nil) {
		go s.readAck(conn)
	}

	for {
		// 获取请求消息
		_, msg, err := conn.ReadMessage()
		if err != nil {
			s.Errorf("server read msg err : %v", err)
			//关闭连接对象
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

		// todo: 给客户端回复一个ack

		//根据消息进行处理
		if s.isAck(&message) {
			s.Infof("conn message read ack msg %v", message)
			conn.appendMsgMq(&message)
		} else {
			conn.message <- &message
		}
	}

}

func (s *Server) isAck(message *Message) bool {
	if message == nil {
		return s.opt.ack != NoAck
	}
	return s.opt.ack != NoAck && message.FrameType != FrameNoAck
}

// 读取消息的ack
func (s *Server) readAck(conn *Conn) {
	for {
		select {
		case <-conn.done:
			s.Infof("close message ack uid %v", conn.Uid)
			return
		default:

		}
		//从队列中读取新的消息
		conn.messageMu.Lock()
		if len(conn.message) == 0 {
			conn.messageMu.Unlock()
			//增加睡眠
			time.Sleep(100 * time.Millisecond)
			continue
		}
		//读取第一条
		message := conn.readMessage[0]

		//判断ack的方式
		switch s.opt.ack {
		case OnlyAck:
			//直接给客户端回复
			s.Send(&Message{
				FrameType: FrameAck,
				Id:        message.Id,
				AckSeq:    message.AckSeq + 1,
			}, conn)
			//进行业务处理
			//把消息从队列中移除
			conn.readMessage = conn.readMessage[1:]
			conn.messageMu.Unlock()
			conn.message <- message
		case RigorAck:
			//先回
			if message.AckSeq == 0 {
				//还未确认
				conn.readMessage[0].AckSeq++
				conn.readMessage[0].ackTime = time.Now()
				s.Send(&Message{
					FrameType: FrameAck,
					Id:        message.Id,
					AckSeq:    message.AckSeq,
				}, conn)
				s.Infof("message ack RigorAck send mid %v ,seq %v , time %v",
					message.Id, message.AckSeq, message.ackTime)
				conn.messageMu.Unlock()
				continue
			}
			//再验证

			//1.客户端返回结果，再一次确认
			//得到客户端的序号
			msgSeq := conn.readMessageSeq[message.Id]
			if msgSeq.AckSeq > message.AckSeq {
				//确认
				conn.readMessage = conn.readMessage[1:]
				conn.messageMu.Unlock()
				conn.message <- message
				s.Infof("message ack RigorAck success mid %v ", message.Id)
				continue
			}
			//2.客户端没有确认,考虑是否超过了ack的确认时间
			val := s.opt.ackTimeout - time.Since(message.ackTime)
			if !message.ackTime.IsZero() && val <= 0 {
				// 2.2 超过结束确认
				delete(conn.readMessageSeq, message.Id)
				conn.readMessage = conn.readMessage[1:]
				conn.messageMu.Unlock()
				continue
			}
			// 2.1 未超过，重新发送
			conn.messageMu.Unlock()
			s.Send(&Message{
				FrameType: FrameAck,
				Id:        message.Id,
				AckSeq:    message.AckSeq,
			}, conn)
			//再睡眠一定的时间
			time.Sleep(3 * time.Second)
		}
	}
}

// 任务的处理
func (s *Server) handlerWrite(conn *Conn) {
	for {
		select {
		case <-conn.done:
			//连接关闭
			return
		case message := <-conn.message:
			switch message.FrameType {
			case FramePing:
				s.Send(&Message{
					FrameType: FramePing,
				}, conn)
			case FrameData:
				//根据请求的method 分发路由并执行
				if handler, ok := s.routes[message.Method]; ok {
					handler(s, conn, message)
				} else {
					//conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("不存在执行的方法 %v 请检查", message.Method)))
					s.Send(&Message{
						FrameType: FrameData,
						Data:      fmt.Sprintf("不存在执行的方法 %v 请检查", message.Method),
					}, conn)
				}
			}
			if s.isAck(message) {
				conn.messageMu.Lock()
				delete(conn.readMessageSeq, message.Id)
				conn.messageMu.Unlock()
			}
		}
	}
}

// AddRouters 添加实际具体要执行的路由
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
