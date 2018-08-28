package log

import (
	"fmt"
	"os"
	"strings"

	"github.com/lestrrat/go-file-rotatelogs"
	// "github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	formatter "github.com/x-cray/logrus-prefixed-formatter"
)

type Log struct {
	log *logrus.Logger
}

func (l *Log) initLocalLog(name string, opts ...option) error {
	caller = findCaller(opts...)

	level, err := findLevel(opts...)
	if nil != err {
		return err
	}

	if findTerminal(opts...) {
		return terminal(level)
	}

	var dir string = ""
	var path string = ""
	if "" == name {
		dir = findLogName(opts...)
		if err := createPath(dir); nil != err {
			return err
		}
		path = fmt.Sprintf("%s_%d", dir, os.Getpid())
		path += ".%Y_%m_%d"
	} else {
		pos := strings.LastIndex(name, "/")
		if -1 == pos {
			dir = "./"
		} else {
			dir = name[:pos+1]
		}
		path = name
	}

	writer, err := rotatelogs.New(
		path,
		rotatelogs.WithMaxAge(findMaxAge(opts...)),
		rotatelogs.WithRotationTime(findRotationTime(opts...)),
		rotatelogs.WithRotationCount(findRotationCount(opts...)),
	)
	if err != nil {
		return err
	}

	if findWatcherEnable(opts...) {
		go watcher(dir, opts...)
	}

	l.log.Out = writer
	l.log.SetLevel(level)
	l.log.Formatter = &formatter.TextFormatter{
		TimestampFormat:  "2006-01-02 15:04:05.00000000",
		ForceColors:      true,
		QuoteEmptyFields: true,
		FullTimestamp:    true,
	}

	return nil
}

func (l *Log) initRemoteLog(opts ...option) error {
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

	l.log.Formatter = &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05.00000000",
	}

	l.log.SetLevel(level)

	if findFluent(opts...) {
		return fluent(l.log, addr, opts...)
	}

	go handle(ctx, l.log, addr, findRemoteProtocolType(opts...))

	return nil
}

func (l *Log) output() {

}

func (l *Log) NewEntry(name string) *logrus.Entry {
	return l.log.WithField("model", name)
}
