package constants

type Mtype int

const (
	TextMtype Mtype = iota
)

type ChatType int

const (
	GroupChatType ChatType = iota + 1
	SingleChatType
)
