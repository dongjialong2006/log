package log

import (
	"context"
	"time"
)

func WithLogName(name string) option {
	return new(logName, name)
}

func WithTerminal(out bool) option {
	return new(logTerminal, out)
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

func WithWatchEnable(enable bool) option {
	return new(logWatcherEnable, enable)
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

func WithFluent(value bool) option {
	return new(logFluent, value)
}
