package logger_test

import (
	logger "immortality/service/logger"
	"os"
	"testing"
)

func TestStdoutLogger(t *testing.T) {
	l := logger.NewStdoutLogger(os.Stdout, "hehe")
	l.Info("haha", "info")
	l.Warn("haha", "warn")
	l.Error("haha", "error")
}
