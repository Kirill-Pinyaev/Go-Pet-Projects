package logger

import (
	"encoding/json"
	"os"
	"time"
)

type ConsoleLogger struct{}

func (ConsoleLogger) log(level, msg string) {
	obj := map[string]string{
		"timestamp": time.Now().Format(time.RFC3339),
		"level":     level,
		"msg":       msg,
	}
	_ = json.NewEncoder(os.Stdout).Encode(obj) // одна строка JSON
}

func (c ConsoleLogger) Info(msg string) {
	c.log("info", msg)
}

func (c ConsoleLogger) Error(err error) {
	c.log("error", err.Error())
}
