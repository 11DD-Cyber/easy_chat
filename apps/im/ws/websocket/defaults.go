package websocket

import (
	"math"
	"time"
)

const (
	infinity                 = time.Duration(math.MaxInt64)
	defaultMaxConnectionIdle = infinity
)

const (
	defaultAckTimeout = 30 * time.Second
)
