package util

import (
	"context"
	"sync"
)

var (
	shutdowncontext *ShutdownContext
	onceLock        sync.Once
)

type ShutdownContext struct {
	Context context.Context
	Cancel  context.CancelFunc
}

func initialize() {
	onceLock.Do(func() {
		ctx, cancel := context.WithCancel(context.Background())
		shutdowncontext = &ShutdownContext{ctx, cancel}
	})
}

func GetShutDownCtx() *ShutdownContext {
	initialize()
	return shutdowncontext
}
