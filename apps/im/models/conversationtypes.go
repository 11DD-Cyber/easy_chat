package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Conversation struct {
	ID             bson.ObjectID `bson:"_id,omitempty"`  // MongoDB 集合主键
	ConversationId string        `bson:"conversationId"` // 会话唯一标识（业务层ID）
	ChatType       int32         `bson:"chatType"`       // 1=单聊，2=群聊
	TargetId       string        `bson:"targetId"`       // 补充：对方/群ID（IM 必需）
	IsShow         bool          `bson:"isShow"`         // 是否显示在会话列表
	Total          int32         `bson:"total"`          // 会话总消息数
	Seq            int64         `bson:"seq"`            // 消息最新序列号
	ToRead         int32         `bson:"toRead"`         // 补充：未读消息数（IM 必需）
	Read           int32         `bson:"read"`           // 补充：已读消息数
	Msg            *ChatLog      `bson:"msg"`            // 会话最后一条消息
	UpdateAt       time.Time     `bson:"updateAt"`       // 会话更新时间
	CreateAt       time.Time     `bson:"createAt"`       // 会话创建时间
}
