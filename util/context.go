package util

import (
	"context"
	"time"
)

func InitConText(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}
