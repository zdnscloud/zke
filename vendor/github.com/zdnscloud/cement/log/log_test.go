package log

import (
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	logger, _ := NewLog4jLogger("test.log", Warn, 0, 0)
	logger.Info("info message")
	logger.Debug("debug message")
	logger.Warn("warn message")
	logger.Error("error message")
	<-time.After(1 * time.Second)
	logger.Close()
}

func TestTermLogger(t *testing.T) {
	InitLogger(Debug)
	defer CloseLogger()
	Fatalf("good boy")
}
