package websocket

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Conn struct {
	idleMu sync.Mutex //空闲操作锁

	Uid string

	*websocket.Conn
	s *Server

	idle              time.Time     //空闲时间
	maxConnectionIdle time.Duration //最大空闲时间

	messageMu      sync.Mutex
	readMessage    []*Message          //读消息的队列，用于ack机制
	readMessageSeq map[string]*Message //消息的序列化

	message chan *Message //通道，ack确认之后将消息发送给任务处理

	done chan struct{}
}

// NewConn 通过http的方式获取连接对象
func NewConn(s *Server, w http.ResponseWriter, r *http.Request) *Conn {
	c, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.Errorf("upgrade err %v", err)
		return nil
	}

	conn := &Conn{
		Conn:              c,                    //连接对象
		s:                 s,                    //服务
		idle:              time.Now(),           //空闲时间
		maxConnectionIdle: s.opt.maxConnectIdle, //最大空闲时间

		readMessage:    make([]*Message, 0, 2),       //队列
		readMessageSeq: make(map[string]*Message, 2), //消息队列的序列化
		message:        make(chan *Message, 1),       //消息的通道

		done: make(chan struct{}), //结束服务的通道
	}

	// 执行心跳检测
	go conn.keepalive()

	return conn
}

// 读取队列，发送消息
func (c *Conn) appendMsgMq(msg *Message) {
	c.messageMu.Lock()
	defer c.messageMu.Unlock()

	//读队列中
	if m, ok := c.readMessageSeq[msg.Id]; ok {
		//已经有消息的记录，该消息已经有ack的确认
		if len(c.readMessage) == 0 {
			//队列中没有该消息
			return
		}

		//msg.AckSeq > m.AckSeq
		if m.AckSeq >= msg.AckSeq {
			//没有进行ack的确认,重复
			return
		}
		c.readMessageSeq[msg.Id] = msg
		return
	}
	// 还没有进行ack的确认,避免客户端重复发送多余的ack消息
	if msg.FrameType == FrameAck {
		return
	}

	//记录好消息，消息的序列号
	c.readMessage = append(c.readMessage, msg)
	c.readMessageSeq[msg.Id] = msg

}

// ReadMessage 连接对象读操作
func (c *Conn) ReadMessage() (messageType int, data []byte, err error) {
	messageType, data, err = c.Conn.ReadMessage()

	c.idleMu.Lock()
	defer c.idleMu.Unlock()
	c.idle = time.Time{}

	return
}

// WriteMessage 连接对象写操作
func (c *Conn) WriteMessage(messageType int, data []byte) error {
	c.idleMu.Lock()
	defer c.idleMu.Unlock()
	// 方法是并不安全
	err := c.Conn.WriteMessage(messageType, data)
	c.idle = time.Now()
	return err
}

func (c *Conn) Close() error {
	select {
	case <-c.done:
	default:
		close(c.done)
	}

	return c.Conn.Close()
}

// 长连接检测机制
func (c *Conn) keepalive() {
	// 定时器
	idleTimer := time.NewTimer(c.maxConnectionIdle)
	defer func() {
		idleTimer.Stop()
	}()

	for {
		select {
		case <-idleTimer.C:
			c.idleMu.Lock()
			idle := c.idle
			if idle.IsZero() { // The connection is non-idle. idle连接处于非空闲状态
				c.idleMu.Unlock()
				idleTimer.Reset(c.maxConnectionIdle)
				continue
			}
			val := c.maxConnectionIdle - time.Since(idle)
			c.idleMu.Unlock()
			if val <= 0 {
				// 超时了
				// The connection has been idle for a duration of keepalive.MaxConnectionIdle or more.
				// Gracefully close the connection.
				c.s.Close(c)
				return
			}
			idleTimer.Reset(val)
		case <-c.done:
			fmt.Println("客户端结束连接")
			return
		}
	}
}
