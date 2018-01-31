package channel

import (
	"fmt"
	"time"
)

const (
	prefix = "channel:"

	defaultCheckInterval = time.Millisecond * 100
	defaultTimeoutDelta = time.Minute
)

var (
	ErrTimeout = fmt.Errorf(prefix + "send/recv timeout")
	ErrClosed = fmt.Errorf(prefix + "access a closed channel")
)