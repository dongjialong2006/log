package log

import (
	"fmt"
	"log/file"
	"log/logger"
	"log/remote"
	"log/types"
	"os"

	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	formatter "github.com/x-cray/logrus-prefixed-formatter"
)

func init() {
	InitLocalLogSystem(WithLogLevel("debug"),
		WithMaxAge(types.DEFAULT_MAX_AGE),
		WithRotationCount(types.DEFAULT_ROTATION_COUNT),
		WithRotationTime(types.DEFAULT_ROTATION_TIME),
		WithWatchEnable(),
		WithCaller(),
	)

	/*
		InitRemoteLogSystem(WithLogLevel("debug"),
			WithCaller(),
			WithRemoteAddr("10.95.135.204:23213"),
			WithRemoteProtocolType(types.TCP),
		)
	*/
}

func New(name string) *Entry {
	e := &Entry{
		Log:    logrus.WithField("model", name),
		Caller: caller,
	}

	for key, field := range fields {
		e.Log = e.Log.WithField(key, field)
	}

	return e
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

	path := findLogName(opts...)
	if err := file.CreatePath(path); nil != err {
		return err
	}

	fields = findFields(opts...)

	path = fmt.Sprintf("%s_%d", path, os.Getpid())

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
		go watcher()
	}

	logrus.SetOutput(&logger.Logger{})

	lfHook := lfshook.NewHook(newWriter(level, writer), &formatter.TextFormatter{
		TimestampFormat:  "2006-01-02 15:04:05.0000",
		ForceColors:      true,
		QuoteEmptyFields: true,
		FullTimestamp:    true,
	})

	logrus.AddHook(lfHook)

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
		TimestampFormat: "2006-01-02 15:04:05.0000",
	})

	logrus.SetLevel(level)

	go remote.Handle(ctx, addr, findRemoteProtocolType(opts...))

	return nil
}
