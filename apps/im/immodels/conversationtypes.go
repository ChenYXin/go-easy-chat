package immodels

import (
	"easy-chat/pkg/constants"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Description: 一个单独会话的数据结构体
type Conversation struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`

	ConversationId string             `bson:"conversationId,omitempty"`
	ChatType       constants.ChatType `bson:"chatType,omitempty"`
	//TargetId       string             `bson:"targetId,omitempty"`
	IsShow bool     `bson:"isShow,omitempty"`
	Total  int      `bson:"total,omitempty"`
	Seq    int64    `bson:"seq"`
	Msg    *ChatLog `bson:"msg,omitempty"`

	UpdateAt time.Time `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
	CreateAt time.Time `bson:"createAt,omitempty" json:"createAt,omitempty"`
}
