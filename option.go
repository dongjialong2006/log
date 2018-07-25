package log

import (
	"context"
	"time"
)

type option interface {
	Get(key string) interface{}
}

type Option struct {
	key   string
	value interface{}
}

func new(key string, value interface{}) *Option {
	return &Option{
		key:   key,
		value: value,
	}
}

func (opt Option) Get(key string) interface{} {
	if key == opt.key {
		return opt.value
	}

	return nil
}

const (
	logAge           = "age"
	logName          = "name"
	logTime          = "time"
	logLevel         = "level"
	logCount         = "count"
	logFields        = "fields"
	logContext       = "context"
	logCaller        = "caller"
	logTerminal      = "terminal"
	logRemoteAddr    = "addr"
	logRemoteProType = "protocol"
	logWatcherEnable = "enable"
	logWatcherByNum  = "watcherByNum"
	logWatcherBySize = "watcherBySize"
)

func WithLogName(name string) option {
	return new(logName, name)
}

func WithTerminal() option {
	return new(logTerminal, true)
}

func WithRemoteAddr(addr string) option {
	return new(logRemoteAddr, addr)
}

func WithRemoteProtocolType(protocol string) option {
	return new(logRemoteProType, protocol)
}

func WithCaller() option {
	return new(logCaller, true)
}

func WithContext(ctx context.Context) option {
	return new(logContext, ctx)
}

func WithWatchEnable() option {
	return new(logWatcherEnable, true)
}

func WithWatchLogsByNum(num int) option {
	return new(logWatcherByNum, num)
}

func WithWatchLogsBySize(size int64) option {
	return new(logWatcherBySize, size)
}

func WithRotationTime(interval time.Duration) option {
	return new(logTime, interval)
}

func WithRotationCount(num int) option {
	return new(logCount, num)
}

func WithMaxAge(interval time.Duration) option {
	return new(logAge, interval)
}

func WithLogLevel(level string) option {
	return new(logLevel, level)
}

func WithFields(fields map[string]interface{}) option {
	return new(logFields, fields)
}
