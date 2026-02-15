package mq

import "easy_chat/pkg/constants"

type MsgChatTransfer struct {
	ConversationId string          `json:"conversationId"`
	SendId         string          `json:"sendId"`
	RecvId         string          `json:"recvId"`
	SendTime       int64           `json:"sendTime"`
	Content        string          `json:"content"`
	Msgtype        constants.MType `json:"mType"`
}
