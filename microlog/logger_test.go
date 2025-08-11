package logger_test

import (
	"errors"
	"microlog/logger"
	"strings"
	"testing"
)

func TestBufferLogger(t *testing.T) {
	bl := &logger.BufferLogger{}
	bl.Error(errors.New("boom"))

	got := bl.Dump()
	wantSubstr := `"level":"error"`
	if !strings.Contains(got, wantSubstr) {
		t.Fatalf("expect %q to contain %q", got, wantSubstr)
	}
}

func TestTypeAssertion(t *testing.T) {
	bl := &logger.BufferLogger{}
	var l logger.Logger = bl

	// пишем хотя бы одно сообщение
	l.Info("hello world")

	if got := bl.Dump(); strings.TrimSpace(got) == "" {
		t.Fatalf("Dump() вернул пустую строку, ожидали JSON-запись")
	}
}
