package gor

import (
	"runtime/debug"

	"github.com/chenjie4255/tools/log"
)

var logger *log.Logger

func init() {
	logger = log.NewLoggerWithSentry("gor_recover")
}

// RunWithRecover 在新的协程中运行FUNC，带recover保护
func RunWithRecover(fn func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logger.WithFields(log.Fields{
					"panic": err,
					"stack": string(debug.Stack()),
				}).Error("recover panic from goroutine!")
			}
		}()

		fn()
	}()
}
