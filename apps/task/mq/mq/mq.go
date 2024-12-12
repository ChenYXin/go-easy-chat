package mq

import "easy-chat/pkg/constants"

// MsgChatTransfer 消息转化的格式，用于消息队列
type MsgChatTransfer struct {
	ConversationId     string `json:"conversationId"`
	constants.ChatType `json:"chatType"`
	SendId             string `json:"sendId"`
	RecvId             string `json:"recvId"`
	SendTime           int64  `json:"sendTime"`

	constants.Mtype `json:"mType"`
	Content         string `json:"content"`
}
