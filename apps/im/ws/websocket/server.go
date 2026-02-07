package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

type Server struct {
	// 路由规则表
	sync.RWMutex
	authentication Authentication
	routes         map[string]HandlerFunc
	addr           string
	connToUser     map[*Conn]string
	userToConn     map[string]*Conn
	upgrader       websocket.Upgrader
	Logger         logx.Logger
}

func (s *Server) AddRoutes(rs []Route) {
	s.Lock()
	defer s.Unlock()
	for _, r := range rs {

		s.routes[r.Method] = r.Handler
	}
}
func NewServer(addr string, auth Authentication) *Server {
	return &Server{
		routes:         make(map[string]HandlerFunc),
		connToUser:     make(map[*Conn]string),
		authentication: auth,
		userToConn:     make(map[string]*Conn),
		addr:           addr,
		upgrader:       websocket.Upgrader{},
		// 初始化日志：绑定上下文
		Logger: logx.WithContext(context.Background()),
	}
}

func (s *Server) ServerWs(w http.ResponseWriter, r *http.Request) {
	// 异常恢复：防止单个WS连接处理panic导致整个服务崩溃
	defer func() {
		if r := recover(); r != nil {
			s.Logger.Errorf("server handler ws recover err %v", r)
		}
	}()
	//鉴权
	if !s.authentication.Auth(w, r) {
		s.Logger.Errorf("authentication failed")
		http.Error(w, "authentication failed", http.StatusUnauthorized)
		return
	}
	// 把HTTP请求升级为WebSocket连接
	conn := NewConn(w, r, s)
	if conn == nil {
		return
	}
	//心跳
	go conn.keepalive()
	//记录连接
	s.addConn(conn, r)
	//读取信息，完成请求
	go s.handlerConn(conn)

}

// 根据连接对象执行任务处理
func (s *Server) handlerConn(conn *Conn) {
	defer func() {
		s.Close(conn)

	}()
	for {
		//获取请求消息
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
		//根据请求消息类型分类处理
		switch message.FrameType {
		case FramePing:
			//心跳包
			s.Send(&Message{FrameType: FramePing}, conn.Conn)
		case FrameData:
			//数据包，继续处理
			if handler, ok := s.routes[message.Method]; ok {
				handler(s, conn, &message)
			} else {
				s.Send(NewErrMessage(fmt.Errorf("method %s not found", message.Method)), conn.Conn)
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
	defer s.RWMutex.Unlock()
	//如果原来就有连接
	if c := s.userToConn[uid]; c != nil {
		//关闭原来的连接
		c.Close()
	}
	s.connToUser[conn] = uid
	s.userToConn[uid] = conn
	s.Logger.Infof("user %s connect, remote: %s", uid, req.RemoteAddr)
}

func (s *Server) GetConns(uids ...string) []*websocket.Conn {
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()
	var res []*websocket.Conn
	if len(uids) == 0 {
		//获取全部
		res = make([]*websocket.Conn, 0, len(s.userToConn))
		for _, conn := range s.userToConn {
			res = append(res, conn.Conn)
		}
	} else {
		//获取部分
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
		//获取部分
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
	//先关闭连接（done通道和底层WS连接）
	conn.Close()
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()
	uid := s.connToUser[conn]
	delete(s.connToUser, conn)
	delete(s.userToConn, uid)
}

func (s *Server) Start() {
	http.HandleFunc("/ws", s.ServerWs)
	http.ListenAndServe(s.addr, nil)
}

func (s *Server) Stop() {

	//关闭所有连接
	s.Lock()
	defer s.Unlock()
	for conn := range s.connToUser {
		conn.Close()
	}
	//清空映射
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
