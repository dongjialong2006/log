package log

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"
	"sync/atomic"
	"time"

	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/sirupsen/logrus"
	formatter "github.com/x-cray/logrus-prefixed-formatter"
)

type Log struct {
	log   *logrus.Logger
	name  string
	index int32
}

func (l *Log) initLocalLog(name string, opts ...option) error {
	caller = findCaller(opts...)

	level, _ := findLevel(opts...)

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

	writer.Write([]byte(""))

	if findWatcherEnable(opts...) {
		go l.watcher(dir, opts...)
	}

	l.log.Out = writer
	l.log.SetLevel(level)
	l.log.Formatter = &formatter.TextFormatter{
		TimestampFormat:  "2006-01-02 15:04:05.00000000",
		ForceColors:      true,
		QuoteEmptyFields: true,
		FullTimestamp:    true,
	}

	l.log.Debugf("log system init success.")
	l.name = writer.CurrentFileName()

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

func (l *Log) watcher(dir string, opts ...option) {
	ctx := findContext(opts...)
	num := findWatchLogsByNum(opts...)
	size := findWatchLogsBySize(opts...)

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
			l.log.Warnf("watcher path:%s is closed.", dir)
			return
		case <-tick:
			files, err := ioutil.ReadDir(dir)
			if err != nil {
				continue
			}

			l.delBySize(size, files)
			l.delByNum(num, dir, files)
		}
	}

	return
}

func (l *Log) delBySize(basic int64, files []os.FileInfo) {
	for _, f := range files {
		if f.IsDir() || f.Name() != l.name {
			continue
		}
		if f.Size() > basic {
			os.Rename(l.name, fmt.Sprintf("%s_%d", l.name, atomic.AddInt32(&l.index, 1)))
		}
	}
}

func (l *Log) delByNum(num int, dir string, files []os.FileInfo) {
	var logs = make(map[string]string)
	var timestamps []string = nil

	for _, f := range files {
		if f.IsDir() || !strings.Contains(f.Name(), l.name) {
			continue
		}
		timestamps = append(timestamps, f.ModTime().String())
		logs[f.ModTime().String()] = path.Join(dir, f.Name())
	}

	sort.Strings(timestamps)
	for i := 0; i < len(timestamps)-num; i++ {
		os.Remove(logs[timestamps[i]])
		l.log.Debugf("remove file:%s.", logs[timestamps[i]])
	}
}

func (l *Log) NewEntry(name string) *Entry {
	e := &Entry{
		log:    l.log.WithField("model", name),
		caller: caller,
	}

	for key, field := range fields {
		e.log = e.log.WithField(key, field)
	}

	return e
}
