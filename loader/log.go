package loader

import "context"

type Logger interface {
	Error(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Debug(msg string, keysAndValues ...interface{})
}

type noopLogger struct{}

func (noopLogger) Error(msg string, keysAndValues ...interface{}) {}
func (noopLogger) Warn(msg string, keysAndValues ...interface{})  {}
func (noopLogger) Info(msg string, keysAndValues ...interface{})  {}
func (noopLogger) Debug(msg string, keysAndValues ...interface{}) {}

var noopLog = noopLogger{}

type ctxValue string

var logKey ctxValue = "log"

func SetLogToCtx(ctx context.Context, log Logger) context.Context {
	return context.WithValue(ctx, logKey, log)
}

func logFromCtx(ctx context.Context) Logger {
	if ctx == nil {
		return noopLog
	}

	if log, ok := ctx.Value(logKey).(Logger); ok {
		return log
	}

	return noopLog
}
