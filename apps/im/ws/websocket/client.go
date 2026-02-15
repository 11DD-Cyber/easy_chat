package websocket

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

type Client interface {
	Close() error
	Send(v any) error
	Read(v any) error
}
type client struct {
	*websocket.Conn
	host    string
	pattern string
	header  http.Header
}

func NewClient(host, pattern string, header http.Header) *client {
	c := &client{
		host:    host,
		pattern: pattern,
		header:  header,
	}
	conn, err := c.dial()
	if err != nil {
		panic(err)
	}
	c.Conn = conn
	return c
}

func (c *client) dial() (*websocket.Conn, error) {
	u := url.URL{Scheme: "ws", Host: c.host, Path: c.pattern}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), c.header)
	return conn, err
}

func (c *client) Send(v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	err = c.WriteMessage(websocket.TextMessage, data)
	if err == nil {
		return nil
	}
	//发送失败了再建立一次连接
	conn, err := c.dial()
	if err != nil {
		return err
	}
	c.Conn = conn
	return c.WriteMessage(websocket.TextMessage, data)
}

func (c *client) Read(v any) error {
	_, msg, err := c.ReadMessage()
	if err != nil {
		return err
	}
	return json.Unmarshal(msg, v)

}
