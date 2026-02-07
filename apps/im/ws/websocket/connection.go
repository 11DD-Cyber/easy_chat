package websocket

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Conn struct {
	*websocket.Conn
	s                 *Server
	mu                sync.Mutex
	idle              time.Time
	maxConnectionIdle time.Duration
	done              chan struct{}
	closed            bool
}

func NewConn(w http.ResponseWriter, r *http.Request, s *Server) *Conn {
	c, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.Logger.Errorf("upgrade connection err %v", err)
		return nil
	}
	conn := &Conn{
		Conn:              c,
		s:                 s,
		idle:              time.Now(),
		maxConnectionIdle: defaultMaxConnectionIdle,
		done:              make(chan struct{}),
		closed:            false,
	}
	return conn
}

func (c *Conn) keepalive() {
	idleTimer := time.NewTimer(c.maxConnectionIdle)
	defer idleTimer.Stop()
	for {
		select {
		case <-idleTimer.C:
			fmt.Printf("idle %v,maxIdle %v \n", time.Since(c.idle), c.maxConnectionIdle)
			val := c.maxConnectionIdle - time.Since(c.idle)
			if val <= 0 {
				c.s.Logger.Infof("connection idle timeout,closing connection")
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

func (c *Conn) ReadMessage() (messageType int, p []byte, err error) {
	messageType, p, err = c.Conn.ReadMessage()
	if err == nil {
		c.mu.Lock()
		c.idle = time.Now() // 更新空闲时间为当前时间
		c.mu.Unlock()
	}
	return
}
func (c *Conn) WriteMessage(messageType int, data []byte) error {
	err := c.Conn.WriteMessage(messageType, data)
	if err == nil {
		c.mu.Lock()
		c.idle = time.Now() // 更新空闲时间为当前时间
		c.mu.Unlock()
	}
	return err
}

// 安全关闭
func (c *Conn) Close() error {
	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return nil
	}
	c.closed = true
	c.mu.Unlock()
	// 关闭done通道，通知keepalive协程退出
	close(c.done)
	// 关闭底层WS连接
	return c.Conn.Close()
}
