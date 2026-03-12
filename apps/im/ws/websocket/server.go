package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

type AckType int

const (
	// 不做 ACK 确认
	NoAck AckType = iota
	// 仅做一次 ACK
	OnlyAck
	// 严格模式，需要完成两次确认
	RigorAck
)

func (t AckType) ToString() string {
	switch t {
	case OnlyAck:
		return "OnlyAck"

	case RigorAck:
		return "RigorAck"

	}
	return "NoAck"
}

type Server struct {
	// 保护用户与连接映射的读写锁
	sync.RWMutex
	authentication Authentication
	routes         map[string]HandlerFunc
	addr           string
	connToUser     map[*Conn]string
	userToConn     map[string]*Conn
	upgrader       websocket.Upgrader
	Logger         logx.Logger
	Ack            AckType
	httpServer     *http.Server
	lifecycle      Lifecycle
}

type Lifecycle interface {
	OnConnect(s *Server, conn *Conn, uid string)
	OnDisconnect(s *Server, conn *Conn, uid string)
}

func (s *Server) AddRoutes(rs []Route) {
	s.Lock()
	defer s.Unlock()
	for _, r := range rs {

		s.routes[r.Method] = r.Handler
	}
}

func (s *Server) SetLifecycle(l Lifecycle) {
	s.lifecycle = l
}
func NewServer(addr string, auth Authentication) *Server {
	return &Server{
		routes:         make(map[string]HandlerFunc),
		connToUser:     make(map[*Conn]string),
		authentication: auth,
		userToConn:     make(map[string]*Conn),
		addr:           addr,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Allow frontend running on different origin to establish websocket connection.
				return true
			},
		},
		// 初始化日志上下文
		Logger: logx.WithContext(context.Background()),
		Ack:    NoAck,
	}
}

func (s *Server) ServerWs(w http.ResponseWriter, r *http.Request) {
	// recover 防护，避免单个连接处理 panic 影响整个服务
	defer func() {
		if r := recover(); r != nil {
			s.Logger.Errorf("server handler ws recover err %v", r)
		}
	}()
	// 鉴权
	if !s.authentication.Auth(w, r) {
		s.Logger.Errorf("authentication failed")
		http.Error(w, "authentication failed", http.StatusUnauthorized)
		return
	}
	// 将 HTTP 请求升级为 WebSocket
	conn := NewConn(w, r, s)
	if conn == nil {
		return
	}

	// 记录连接
	s.addConn(conn, r)
	// 启动读写协程，处理后续消息
	go s.handlerConn(conn)

}

// 基于连接对象，持续读取并处理消息
func (s *Server) handlerConn(conn *Conn) {
	defer func() {
		s.Close(conn)

	}()
	go s.handleWrite(conn)
	if s.Ack != NoAck {
		go s.readAck(conn)
	}
	for {
		// 读取客户端消息
		_, msg, err := conn.ReadMessage()
		if err != nil {
			s.Logger.Errorf("websocket conn read message err %v", err)
			return
		}
		var message Message
		if err := json.Unmarshal(msg, &message); err != nil {
			s.Send(NewErrMessage(err), conn.Conn)
			continue
		}
		switch message.FrameType {
		case FramePing:
			s.Send(&Message{FrameType: FramePing}, conn.Conn)
			continue
		case FrameCAck:
			conn.handleCAck(&message)
			continue
		}
		// 根据帧类型分类处理
		if s.Ack != NoAck && message.FrameType != FrameNoAck {
			// 放入等待 ACK 的队列
			s.Logger.Infof("conn message write in msgMq %v", message)
			conn.appendMsgMq(&message)

		} else {
			conn.message <- &message
		}
	}
}
func (s *Server) readAck(conn *Conn) {
	for {
		select {
		case <-conn.done:
			// 连接已关闭
			s.Logger.Infof("close message ack uid ")
			return
		default:
			conn.messageMu.Lock()
			if len(conn.readMessages) == 0 {
				conn.messageMu.Unlock()
				// 队列为空，等待 100ms
				time.Sleep(100 * time.Microsecond)
				continue

			}
			// 取出队列第一条
			message := conn.readMessages[0]
			// 根据 ACK 策略处理
			switch s.Ack {
			case OnlyAck:
				s.Send(&Message{
					FrameType: FrameAck,
					AckSeq:    message.AckSeq + 1,
					Id:        message.Id,
				}, conn.Conn)
				// 只发送 ACK，稍后交给业务处理
				conn.readMessages = conn.readMessages[1:]
				conn.messageMu.Unlock()
				conn.message <- message
				s.Logger.Infof("message ack OnlyAck send success id %v", message.Id)
			case RigorAck:
				if entry := conn.readMessageSeq[message.Id]; entry != nil {
					if message.AckSeq == 0 {
						message.AckSeq = 1
						entry.AckSeq = 1
						entry.ackTime = time.Now()
						s.Send(&Message{FrameType: FrameAck, AckSeq: 1, Id: message.Id}, conn.Conn)
						conn.messageMu.Unlock()
						s.Logger.Infof("message ack RigorAck send")
						continue
					}
					if entry.ackConfirmed == true { // 客户端已完成最终确认
						conn.readMessages = conn.readMessages[1:]
						conn.messageMu.Unlock()
						conn.message <- message
						s.Logger.Infof("message ack RigorAck confirmed id %v", message.Id)
						continue
					}
					if time.Since(entry.ackTime) > defaultAckTimeout {
						entry.ackTime = time.Now()
						s.Send(&Message{FrameType: FrameAck, AckSeq: message.AckSeq, Id: message.Id}, conn.Conn)

					}
					conn.messageMu.Unlock()
				}
				time.Sleep(time.Millisecond * 20)
			}
		}
	}
}
func (s *Server) handleWrite(conn *Conn) {
	for {
		select {
		case <-conn.done:
			// 连接结束
			return
		case message := <-conn.message:
			// 按帧类型下发
			switch message.FrameType {
			case FramePing:
				// 直接回复 Pong
				s.Send(&Message{FrameType: FramePing}, conn.Conn)
			case FrameData, FrameNoAck:
				// 业务数据
				if handler, ok := s.routes[message.Method]; ok {
					handler(s, conn, message)
				} else {
					s.Send(&Message{
						FrameType: FrameData,
						Data:      fmt.Sprintf("未找到路由: %v", message.Method),
					}, conn.Conn)
				}
			}
			if s.Ack != NoAck {
				// 删除已完成 ACK 的序列
				conn.messageMu.Lock()
				delete(conn.readMessageSeq, message.Id)
				conn.messageMu.Unlock()

			}
		}
	}
}
func (s *Server) addConn(conn *Conn, req *http.Request) {
	uid := s.authentication.UserId(req)
	if uid == "" {
		s.Logger.Errorf("empty user id for conn %v", conn)
		return
	}
	s.RWMutex.Lock()
	// 若已有连接，关闭旧连接
	if c := s.userToConn[uid]; c != nil {
		c.Close()
	}
	s.connToUser[conn] = uid
	s.userToConn[uid] = conn
	s.RWMutex.Unlock()
	s.Logger.Infof("user %s connect, remote: %s", uid, req.RemoteAddr)
	if s.lifecycle != nil {
		go s.lifecycle.OnConnect(s, conn, uid)
	}
}

func (s *Server) GetConns(uids ...string) []*websocket.Conn {
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()
	var res []*websocket.Conn
	if len(uids) == 0 {
		// 返回全部连接
		res = make([]*websocket.Conn, 0, len(s.userToConn))
		for _, conn := range s.userToConn {
			res = append(res, conn.Conn)
		}
	} else {
		// 仅返回指定用户连接
		res = make([]*websocket.Conn, 0, len(uids))
		for _, uid := range uids {
			if conn, ok := s.userToConn[uid]; ok {
				res = append(res, conn.Conn)
			}
		}
	}
	return res
}

func (s *Server) GetUsers(conns ...*Conn) []string {
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()
	var res []string
	if len(conns) == 0 {
		res = make([]string, 0, len(s.connToUser))
		for _, uid := range s.connToUser {
			res = append(res, uid)
		}
	} else {
		// 获取部分连接对应的用户
		res = make([]string, 0, len(conns))
		for _, conn := range conns {
			if uid, ok := s.connToUser[conn]; ok {
				res = append(res, uid)
			}
		}
	}
	return res
}

func (s *Server) Close(conn *Conn) {
	// 先关闭 done 通道与底层 WebSocket
	conn.Close()
	s.RWMutex.Lock()
	uid := s.connToUser[conn]
	delete(s.connToUser, conn)
	delete(s.userToConn, uid)
	s.RWMutex.Unlock()
	if uid != "" && s.lifecycle != nil {
		go s.lifecycle.OnDisconnect(s, conn, uid)
	}
}

func (s *Server) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", s.ServerWs)
	s.httpServer = &http.Server{
		Addr:    s.addr,
		Handler: mux,
	}
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.Logger.Errorf("websocket server listen error: %v", err)
	}
}

func (s *Server) Stop() {

	// 关闭全部连接
	s.Lock()
	defer s.Unlock()
	for conn := range s.connToUser {
		conn.Close()
	}
	// 清空映射
	s.connToUser = make(map[*Conn]string)
	s.userToConn = make(map[string]*Conn)
	fmt.Println("stop server")
}

func (s *Server) SendByUserId(msg interface{}, sendIds ...string) error {
	if len(sendIds) == 0 {
		return nil
	}
	return s.Send(msg, s.GetConns(sendIds...)...)
}
func (s *Server) Send(msg interface{}, conns ...*websocket.Conn) error {
	if len(conns) == 0 {
		return nil
	}
	data, err := json.Marshal(msg)
	if err != nil {
		s.Logger.Errorf("marshal msg err: %v", err)
		return err
	}
	for _, conn := range conns {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			s.Logger.Errorf("send message err %v", err)
			continue
		}
	}
	return nil
}
