package websocket

type FrameType uint8

const (
	FrameData FrameType = 0x0
	FramePing FrameType = 0x1
	FrameErr  FrameType = 0x2
)

type Message struct {
	FrameType `json:"frameType"`
	Method    string      `json:"method"`
	FromId    string      `json:"fromId"`
	Data      interface{} `json:"data"`
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
