package log

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/sirupsen/logrus"
	formatter "github.com/x-cray/logrus-prefixed-formatter"
)

type Log struct {
	log   *logrus.Logger
	paths map[string]int
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

	var name string = ""
	pos := strings.LastIndex(dir, "/")
	if -1 != pos {
		name = dir[pos+1:]
		dir = dir[:pos+1]
	} else {
		name = dir
		dir = "./"
	}

	l.log.Debugf("log name:%s, path:%s.", name, dir)

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

			l.delBySize(name, size, num, dir, files)
			l.delByNum(name, num, dir, files)
		}
	}

	return
}

func (l *Log) delBySize(name string, basic int64, num int, dir string, files []os.FileInfo) {
	var logs = make(map[string]string)
	var timestamps []string = nil
	var size int64 = 0

	for _, f := range files {
		if f.IsDir() || !strings.Contains(f.Name(), name) {
			continue
		}
		size += f.Size()
		timestamps = append(timestamps, f.ModTime().String())
		logs[f.ModTime().String()] = path.Join(dir, f.Name())
	}

	sort.Strings(timestamps)
	if size >= basic {
		for i := 0; i < len(timestamps); i++ {
			if len(timestamps) == 1 {
				l.replace(logs[timestamps[i]])
				break
			}
			os.Remove(logs[timestamps[i]])
			l.log.Debugf("remove file:%s.", logs[timestamps[i]])
			if len(timestamps) < num {
				break
			}
		}
	}
}

func (l *Log) replace(name string) {
	num, ok := l.paths[name]
	if !ok {
		num = 1
	} else {
		num++
	}
	l.paths[name] = num

	os.Rename(name, fmt.Sprintf("%s_%d", name, num))
}

func (l *Log) delByNum(name string, num int, dir string, files []os.FileInfo) {
	var logs = make(map[string]string)
	var timestamps []string = nil

	for _, f := range files {
		if f.IsDir() || !strings.Contains(f.Name(), name) {
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
