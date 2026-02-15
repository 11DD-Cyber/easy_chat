package websocket

import "time"

type FrameType uint8

const (
	FrameData  FrameType = 0x0
	FramePing  FrameType = 0x1
	FrameErr   FrameType = 0x2
	FrameAck   FrameType = 0x3
	FrameNoAck FrameType = 0x4
	FrameCAck  FrameType = 0x5
)

type Message struct {
	FrameType `json:"frameType"`
	Id        string      `json:"id"`
	Method    string      `json:"method"`
	FromId    string      `json:"fromId"`
	Data      interface{} `json:"data"`
	AckSeq    int         `json:"ackSeq,omitempty"`
	ackTime   time.Time   `json:"-"`
	errCount  int         `json:"-"`
	ackConfirmed bool
}

func NewMessage(formId string, data interface{}) *Message {
	return &Message{
		FromId:    formId,
		FrameType: FrameData,
		Data:      data,
	}
}
func NewErrMessage(err error) *Message {
	return &Message{
		FrameType: FrameErr,
		Data:      err.Error(),
	}
}
