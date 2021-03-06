package log

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/evalphobia/logrus_fluent"
	"github.com/sirupsen/logrus"
	formatter "github.com/x-cray/logrus-prefixed-formatter"
)

var caller bool = false
var fields map[string]interface{} = nil

func fluent(log *logrus.Logger, addr string, opt ...option) error {
	pos := strings.Index(addr, ":")
	if -1 == pos {
		return fmt.Errorf("addr format error.")
	}

	port, err := strconv.Atoi(addr[pos+1:])
	if nil != err {
		return err
	}

	hook, err := logrus_fluent.NewWithConfig(logrus_fluent.Config{
		Host: addr[:pos],
		Port: port,
	})
	if err != nil {
		return err
	}

	hook.SetLevels([]logrus.Level{
		logrus.PanicLevel,
		logrus.ErrorLevel,
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.DebugLevel,
		logrus.FatalLevel,
	})

	name := findLogName(opt...)
	if "" == name {
		return fmt.Errorf("tag is empty.")
	}

	pos = strings.LastIndex(name, "/")
	if -1 != pos {
		name = name[pos+1:]
	}

	hook.SetTag(name)
	hook.AddFilter("error", logrus_fluent.FilterError)

	if nil == log {
		logrus.AddHook(hook)
	} else {
		log.AddHook(hook)
	}

	return nil
}

func terminal(level logrus.Level) error {
	logrus.SetLevel(level)
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&formatter.TextFormatter{
		ForceFormatting:  true,
		FullTimestamp:    true,
		QuoteEmptyFields: true,
	})

	return nil
}

func findLevel(opts ...option) logrus.Level {
	for _, opt := range opts {
		if nil == opt {
			continue
		}

		if value := opt.Get(logLevel); nil != value {
			level, _ := logrus.ParseLevel(value.(string))
			return level
		}
	}

	return logrus.InfoLevel
}

func findTerminal(opts ...option) bool {
	for _, opt := range opts {
		if nil == opt {
			continue
		}

		if value := opt.Get(logTerminal); nil != value {
			return value.(bool)
		}
	}

	return true
}

func findLogName(opts ...option) string {
	for _, opt := range opts {
		if nil == opt {
			continue
		}

		if value := opt.Get(logName); nil != value {
			return value.(string)
		}
	}

	return DEFAULT_LOG_NAME
}

func findRotationTime(opts ...option) time.Duration {
	for _, opt := range opts {
		if nil == opt {
			continue
		}

		if value := opt.Get(logTime); nil != value {
			return value.(time.Duration)
		}
	}

	return DEFAULT_ROTATION_TIME
}

func findRotationCount(opts ...option) int {
	for _, opt := range opts {
		if nil == opt {
			continue
		}

		if value := opt.Get(logCount); nil != value {
			return value.(int)
		}
	}

	return DEFAULT_ROTATION_COUNT
}

func findWatcherEnable(opts ...option) bool {
	for _, opt := range opts {
		if nil == opt {
			continue
		}

		if value := opt.Get(logWatcherEnable); nil != value {
			return value.(bool)
		}
	}

	return false
}

func findWatchLogsByNum(opts ...option) int {
	for _, opt := range opts {
		if nil == opt {
			continue
		}

		if value := opt.Get(logWatcherByNum); nil != value {
			return value.(int)
		}
	}

	return DEFAULT_WATCHER_FILES_BY_NUM
}

func findWatchLogsBySize(opts ...option) int64 {
	for _, opt := range opts {
		if nil == opt {
			continue
		}

		if value := opt.Get(logWatcherBySize); nil != value {
			return value.(int64)
		}
	}

	return DEFAULT_WATCHER_FILES_BY_SIZE
}

func findCaller(opts ...option) bool {
	for _, opt := range opts {
		if nil == opt {
			continue
		}

		if value := opt.Get(logCaller); nil != value {
			return value.(bool)
		}
	}

	return false
}

func findContext(opts ...option) context.Context {
	for _, opt := range opts {
		if nil == opt {
			continue
		}

		if value := opt.Get(logContext); nil != value {
			return value.(context.Context)
		}
	}

	return context.Background()
}

func findFields(opts ...option) map[string]interface{} {
	for _, opt := range opts {
		if nil == opt {
			continue
		}

		if value := opt.Get(logFields); nil != value {
			return value.(map[string]interface{})
		}
	}

	return nil
}

func findRemoteAddr(opts ...option) string {
	for _, opt := range opts {
		if nil == opt {
			continue
		}

		if value := opt.Get(logRemoteAddr); nil != value {
			return value.(string)
		}
	}

	return ""
}

func findRemoteProtocolType(opts ...option) string {
	for _, opt := range opts {
		if nil == opt {
			continue
		}

		if value := opt.Get(logRemoteProType); nil != value {
			return value.(string)
		}
	}

	return "tcp"
}

func findFluent(opts ...option) bool {
	for _, opt := range opts {
		if nil == opt {
			continue
		}

		if value := opt.Get(logFluent); nil != value {
			return value.(bool)
		}
	}

	return false
}
