package logger

type Logger interface {
	Info(msg string)
	Error(err error)
}
