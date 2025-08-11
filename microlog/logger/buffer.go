// logger/console.go
package logger

import (
	"bytes"
	"encoding/json"
	"time"
)

type BufferLogger struct {
	buf bytes.Buffer
}

func (b *BufferLogger) log(level, msg string) {
	obj := map[string]string{
		"timestamp": time.Now().Format(time.RFC3339),
		"level":     level,
		"msg":       msg,
	}
	_ = json.NewEncoder(&b.buf).Encode(obj)
}

func (b *BufferLogger) Info(msg string) { b.log("info", msg) }
func (b *BufferLogger) Error(err error) { b.log("error", err.Error()) }

func (b *BufferLogger) Dump() string {
	out := b.buf.String()
	b.buf.Reset()
	return out
}
