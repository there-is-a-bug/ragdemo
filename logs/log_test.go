package logs

import (
	"context"
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	l = &logger{
		f:    nil,
		dir:  "testdata",
		m:    make(chan *Message),
		name: "test",
	}
	ctx, cancel := context.WithCancel(context.Background())
	go l.listen(ctx)
	defer cancel()
	Warn("test %v", "warn")
	Error("test %v", "error")
	Info("test %v", "info")
	time.Sleep(1 * time.Second)
}
