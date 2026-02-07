package ws

import "easy_chat/pkg/constants"

type (
	Msg struct {
		constants.MType `mapstructure:"mType" json:"msgType"`
		Content         string `mapstructure:"content" json:"msgContent"`
	}
)
type (
	Chat struct {
		ConversationId string `mapstructure:"conversationId" json:"conversationId"`
		SendId         string `mapstructure:"sendId" json:"sendId"`
		RecvId         string `mapstructure:"recvId" json:"recvId"`
		Msg            `mapstructure:"msg" json:"msg"`
		SendTime       int64 `mapstructure:"sendTime" json:"sendTime"`
	}
)
