/*
存放kafka的消息格式
*/
package mq

import "easy-chat/pkg/constants"

type MsgChatTransfer struct {
	ConversationId     string            `json:"conversationId"` // 会话id
	constants.ChatType `json:"chatType"` // 聊天类型（私聊，群聊）
	SendId             string            `json:"sendId"` //发送者id
	RecvId             string            `json:"recvId"`
	RecvIds            []string          `json:"recvIds"`
	SendTime           int64             `json:"sendTime"`

	constants.Mtype `json:"mType"` // 消息类型
	Content         string         `json:"content"`
}

type MsgMarkRead struct {
	constants.ChatType `json:"chatType"`
	ConversationId     string   `json:"conversationId"`
	SendId             string   `json:"sendId"` //发送者id
	RecvId             string   `json:"recvId"`
	MsgIds             []string `json:"msgIds"`
}
