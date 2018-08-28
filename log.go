package log

import (
	"fmt"
	"os"

	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/sirupsen/logrus"
	formatter "github.com/x-cray/logrus-prefixed-formatter"
)

const (
	LOCAL  = "local"
	REMOTE = "remote"
)

func New(name string) *Entry {
	e := &Entry{
		log:    logrus.WithField("model", name),
		caller: caller,
	}

	for key, field := range fields {
		e.log = e.log.WithField(key, field)
	}

	return e
}

func NewLog(name string, opts ...option) (*Log, error) {
	log := &Log{
		log: logrus.New(),
	}

	if err := log.init(name, opts...); nil != err {
		return nil, err
	}

	return log, nil
}

func InitLocalLogSystem(opts ...option) error {
	caller = findCaller(opts...)

	level, err := findLevel(opts...)
	if nil != err {
		return err
	}

	if findTerminal(opts...) {
		return terminal(level)
	}

	dir := findLogName(opts...)
	if err := createPath(dir); nil != err {
		return err
	}

	path := fmt.Sprintf("%s_%d", dir, os.Getpid())

	writer, err := rotatelogs.New(
		path+".%Y_%m_%d",
		rotatelogs.WithMaxAge(findMaxAge(opts...)),             // 文件最大保存时间
		rotatelogs.WithRotationTime(findRotationTime(opts...)), // 日志切割时间间隔
		rotatelogs.WithRotationCount(findRotationCount(opts...)),
	)
	if err != nil {
		return err
	}

	if findWatcherEnable(opts...) {
		go watcher(dir, opts...)
	}

	logrus.SetOutput(writer)
	logrus.SetLevel(level)
	logrus.SetFormatter(&formatter.TextFormatter{
		TimestampFormat:  "2006-01-02 15:04:05.00000000",
		ForceColors:      true,
		QuoteEmptyFields: true,
		FullTimestamp:    true,
	})

	return nil
}

func InitRemoteLogSystem(opts ...option) error {
	ctx := findContext(opts...)
	caller = findCaller(opts...)

	level, err := findLevel(opts...)
	if nil != err {
		return err
	}

	if findTerminal(opts...) {
		return terminal(level)
	}

	fields = findFields(opts...)

	addr := findRemoteAddr(opts...)
	if "" == addr {
		return fmt.Errorf("addr is empty.")
	}

	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05.00000000",
	})

	logrus.SetLevel(level)
	logrus.SetOutput(&output{})

	if findFluent(opts...) {
		return fluent(addr, opts...)
	}

	go handle(ctx, addr, findRemoteProtocolType(opts...))

	return nil
}
