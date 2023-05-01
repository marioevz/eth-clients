package utils

import (
	"context"
	"fmt"
	"time"
)

func Shorten(s string) string {
	start := s[0:6]
	end := s[62:]
	return fmt.Sprintf("%s..%s", start, end)
}

type Logging interface {
	Logf(format string, values ...interface{})
}

// Maximum seconds to wait on an RPC request
var rpcTimeoutSeconds = 5

func ContextTimeoutRPC(
	parent context.Context,
) (context.Context, context.CancelFunc) {
	return context.WithTimeout(
		parent,
		time.Second*time.Duration(rpcTimeoutSeconds),
	)
}
