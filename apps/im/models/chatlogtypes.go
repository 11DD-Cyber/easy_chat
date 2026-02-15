package models

import (
	"easy_chat/pkg/constants"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ChatLog struct {
	ID             bson.ObjectID   `bson:"_id" json:"id"`
	ConversationId string          `bson:"conversationId" json:"conversationId"`
	SendId         string          `bson:"sendId" json:"sendId"`
	RecvId         string          `bson:"recvId" json:"recvId"`
	MsgFrom        int             `bson:"msgFrom" json:"msgFrom"`
	MsgType        constants.MType `bson:"msgType" json:"msgType"`
	MsgContent     string          `bson:"msgContent" json:"msgContent"`
	ChatType       int32           `bson:"chatType" json:"chatType,omitempty"`
	SendTime       int64           `bson:"sendTime" json:"sendTime"`
	ReadRecords    []byte          `bson:"readRecords" json:"readRecords,omitempty"`
	Status         int             `bson:"status" json:"status"`
	UpdateAt       time.Time       `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
	CreateAt       time.Time       `bson:"createAt,omitempty" json:"createAt,omitempty"`
}
