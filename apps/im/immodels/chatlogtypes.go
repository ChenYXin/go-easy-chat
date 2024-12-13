package immodels

import (
	"easy-chat/pkg/constants"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 一次性能读多少消息，最大的
var DefaultChatLogLimit int64 = 100

type ChatLog struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`

	ConversationId string             `bson:"conversationId"`
	SendId         string             `bson:"sendId"`
	RecvId         string             `bson:"recvId"`
	MsgFrom        int                `bson:"msgFrom"`
	ChatType       constants.ChatType `bson:"chatType"` //SingleChatType：单聊、GroupChatType：群聊
	MsgType        constants.Mtype    `bson:"msgType"`  //文字
	MsgContent     string             `bson:"msgContent"`
	SendTime       int64              `bson:"sendTime"`
	Status         int                `bson:"status"`
	ReadRecords    []byte             `bson:"readRecords"` // 已读未读

	//Todo : Fill your own fields
	UpdateAt time.Time `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
	CreateAt time.Time `bson:"createAt,omitempty" json:"createAt,omitempty"`
}
