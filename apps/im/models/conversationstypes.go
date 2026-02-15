package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Conversations struct {
	ID               bson.ObjectID            `bson:"_id,omitempty"`
	UserId           string                   `bson:"userId"`           // 关联的用户ID
	ConversationList map[string]*Conversation `bson:"conversationList"` // 会话ID → 会话详情（和 proto 一致）
	UpdateAt         time.Time                `bson:"updateAt"`         // 列表更新时间
	CreateAt         time.Time                `bson:"createAt"`         // 列表创建时间
}
