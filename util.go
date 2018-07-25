package log

import (
	"context"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	formatter "github.com/x-cray/logrus-prefixed-formatter"
)

var log = New("log")
var caller bool = false
var fields map[string]interface{} = nil

type logger struct {
	url      string
	protocol string
}

func (l *logger) Write(p []byte) (n int, err error) {
	if "" == l.url {
		return 0, nil
	}

	if HTTP == l.protocol {
		_, err = post(l.url, string(p))
	}

	if HTTPS == l.protocol {
		_, err = post(l.url, string(p))
	}

	return
}

func newWriter(level logrus.Level, writer *rotatelogs.RotateLogs) lfshook.WriterMap {
	var handles = make(lfshook.WriterMap)
	switch level {
	case 0:
		logrus.SetLevel(logrus.PanicLevel)
		handles[logrus.PanicLevel] = writer
	case 1:
		logrus.SetLevel(logrus.FatalLevel)
		handles[logrus.PanicLevel] = writer
		handles[logrus.FatalLevel] = writer
	case 2:
		logrus.SetLevel(logrus.ErrorLevel)
		handles[logrus.PanicLevel] = writer
		handles[logrus.FatalLevel] = writer
		handles[logrus.ErrorLevel] = writer
	case 3:
		logrus.SetLevel(logrus.WarnLevel)
		handles[logrus.PanicLevel] = writer
		handles[logrus.FatalLevel] = writer
		handles[logrus.ErrorLevel] = writer
		handles[logrus.WarnLevel] = writer
	case 4:
		logrus.SetLevel(logrus.InfoLevel)
		handles[logrus.PanicLevel] = writer
		handles[logrus.FatalLevel] = writer
		handles[logrus.ErrorLevel] = writer
		handles[logrus.WarnLevel] = writer
		handles[logrus.InfoLevel] = writer
	case 5:
		logrus.SetLevel(logrus.DebugLevel)
		handles[logrus.PanicLevel] = writer
		handles[logrus.FatalLevel] = writer
		handles[logrus.ErrorLevel] = writer
		handles[logrus.WarnLevel] = writer
		handles[logrus.InfoLevel] = writer
		handles[logrus.DebugLevel] = writer
	default:
		return nil
	}

	return handles
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

func watcher(opts ...option) {
	ctx := findContext(opts...)
	num := findWatchLogsByNum(opts...)
	size := findWatchLogsBySize(opts...)
	dir := findLogName(opts...)

	pos := strings.LastIndex(dir, "/")
	if -1 != pos {
		dir = dir[:pos+1]
	} else {
		dir = "./"
	}

	tick := time.Tick(time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick:
			files, err := ioutil.ReadDir(dir)
			if err != nil {
				continue
			}

			var logs = make(map[string]string)
			var timestamps []string = nil

			tmp := int64(0)
			for _, f := range files {
				if f.IsDir() || !strings.Contains(f.Name(), "_") || !strings.Contains(f.Name(), ".") {
					continue
				}
				tmp += f.Size()
				timestamps = append(timestamps, f.ModTime().String())
				logs[f.ModTime().String()] = path.Join(dir, f.Name())
			}

			delBySize(size, num, tmp, timestamps, logs)
			delByNum(num, timestamps, logs)
		}
	}

	return
}

func delBySize(basic int64, num int, size int64, timestamps []string, logs map[string]string) {
	sort.Strings(timestamps)
	if size >= basic {
		for i := 0; i < len(timestamps)-2; i++ {
			os.Remove(logs[timestamps[i]])
			log.Debugf("remove file:%s.", logs[timestamps[i]])
			if len(timestamps) < num {
				break
			}
		}
	}
}

func delByNum(num int, timestamps []string, logs map[string]string) {
	sort.Strings(timestamps)
	for i := 0; i < len(timestamps)-num; i++ {
		os.Remove(logs[timestamps[i]])
		log.Debugf("remove file:%s.", logs[timestamps[i]])
	}
}

func findLevel(opts ...option) (logrus.Level, error) {
	for _, opt := range opts {
		if nil == opt {
			continue
		}

		if value := opt.Get(logLevel); nil != value {
			return logrus.ParseLevel(value.(string))
		}
	}

	return logrus.InfoLevel, nil
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

	return false
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

func findMaxAge(opts ...option) time.Duration {
	for _, opt := range opts {
		if nil == opt {
			continue
		}

		if value := opt.Get(logAge); nil != value {
			return value.(time.Duration)
		}
	}

	return DEFAULT_MAX_AGE
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
