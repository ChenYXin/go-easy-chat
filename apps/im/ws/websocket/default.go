package websocket

import (
	"math"
	"time"
)

const (
	defaultMaxConnectionIdle = time.Duration(math.MaxInt64) //默认最大空闲时间
	defaultAckTimeout        = 30 * time.Second             //默认ack超时时间
)
