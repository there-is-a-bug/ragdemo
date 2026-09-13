package logs

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"time"
)

type logger struct {
	f    *os.File
	m    chan *Message
	dir  string
	name string
}

type mtype int

const (
	system mtype = iota
	info
	warn
	erro
	defaultDir = "./log"
)

type Message struct {
	m    string
	t    mtype
	loc  string
	line int
}

func (m mtype) String() string {
	switch m {
	case system:
		return "system"
	case info:
		return "info"
	case warn:
		return "warn"
	case erro:
		return "error"
	default:
		return "unknown"
	}
}

func (m *Message) format() string {
	s := m.t.String()
	t := time.Now().Format("15:04:05")
	return fmt.Sprintf("[%s %d] %s: %s %s", m.loc, m.line, t, s, m.m)
}

var l *logger

func InitLogger(ctx context.Context) {
	l = &logger{}
	l.m = make(chan *Message)
	l.dir = defaultDir
	err := os.MkdirAll(l.dir, 0755)
	if err != nil {
		panic(err)
	}
	err = l.open()
	if err != nil {
		panic(err)
	}
	go l.listen(ctx)
}

func (l *logger) listen(ctx context.Context) {
	for {
		select {
		case msg := <-l.m:
			l.do(msg)
		case <-ctx.Done():
			l.do(&Message{
				m: "logger exit",
				t: system,
			})
			l.close()
			return
		}
	}
}

func (l *logger) do(m *Message) {
	if m == nil {
		return
	}
	//t := time.Now().Format("01-02")
	//if l.name != t {
	//	l.close()
	//	_ = l.open()
	//}
	//
	//if l.f == nil {
	//	err := l.open()
	//	if err != nil {
	//		return
	//	}
	//}
	l.write(m)
}

func (l *logger) open() error {
	name := time.Now().Format("01-02")
	l.name = name
	path := l.dir + "/" + name
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	l.f = file
	return nil
}

func (l *logger) close() {
	if l.f == nil {
		return
	}
	_ = l.f.Sync()
	_ = l.f.Close()
	l.f = nil
}

func Info(m string, args ...any) {
	_, file, line, _ := runtime.Caller(1)
	go func() {
		l.m <- &Message{
			m:    fmt.Sprintf(m, args...),
			t:    info,
			loc:  file,
			line: line,
		}
	}()
}

func Warn(m string, args ...any) {
	_, file, line, _ := runtime.Caller(1)
	go func() {
		l.m <- &Message{
			m:    fmt.Sprintf(m, args...),
			t:    warn,
			loc:  file,
			line: line,
		}
	}()
}

func Error(m string, args ...any) {
	_, file, line, _ := runtime.Caller(1)
	go func() {
		l.m <- &Message{
			m:    fmt.Sprintf(m, args...),
			t:    erro,
			loc:  file,
			line: line,
		}
	}()

}

func (l *logger) write(msg *Message) {
	log := msg.format() + "\n"
	_, _ = l.f.WriteString(log)
	fmt.Println(log)
}
