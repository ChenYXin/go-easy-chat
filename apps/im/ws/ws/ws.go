package ws

import "easy-chat/pkg/constants"

type (
	Msg struct {
		constants.Mtype `mapstructure:"mType"`
		Content         string `mapstructure:"content"`
	}
	Chat struct {
		ConversationId     string `mapstructure:"conversationId"`
		SendId             string `mapstructure:"sendId"`
		RecvId             string `mapstructure:"recvId"`
		SendTime           int64  `mapstructure:"sendTime"`
		constants.ChatType `mapstructure:"chatType"`
		Msg                `mapstructure:"msg"`
	}

	Push struct {
		ConversationId     string `mapstructure:"conversationId"`
		constants.ChatType `mapstructure:"chatType"`
		SendId             string `mapstructure:"sendId"`
		RecvId             string `mapstructure:"recvId"`
		SendTime           int64  `mapstructure:"sendTime"`

		constants.Mtype `mapstructure:"mType"`
		Content         string `mapstructure:"content"`
	}
)
