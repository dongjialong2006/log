package log

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync/atomic"
	"time"

	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	formatter "github.com/x-cray/logrus-prefixed-formatter"
)

const filePathEmpty = "file path is empty."

func (l *Log) logFileName() string {
	if !l.self {
		year, month, day := time.Now().Date()
		return fmt.Sprintf("%s_%d.%d_%02d_%02d", l.formt, os.Getpid(), year, month, day)
	}

	pos := strings.LastIndex(l.formt, "/")
	if -1 == pos {
		return fmt.Sprintf("./%s", l.formt)
	}

	return l.formt
}

func (l *Log) initLocalLogSystem(name string, opts ...option) error {
	level := findLevel(opts...)
	caller = findCaller(opts...)

	if findTerminal(opts...) {
		l.log = logrus.StandardLogger()
		return terminal(level)
	}

	if err := l.defPath(name, level, opts...); err != nil {
		return err
	}

	if findWatcherEnable(opts...) {
		go l.watch(opts...)
	}

	return nil
}

func (l *Log) defPath(name string, level logrus.Level, opts ...option) error {
	if "" != name {
		l.self = true
		l.formt = name
	} else {
		l.formt = findLogName(opts...)
	}

	name = l.logFileName()
	l.name = filepath.Base(name)

	l.path = filepath.Dir(name)

	l.hook = lfshook.NewHook(name, &formatter.TextFormatter{
		TimestampFormat:  "2006-01-02 15:04:05.000000",
		ForceColors:      true,
		QuoteEmptyFields: true,
		FullTimestamp:    true,
	})

	out, err := newOutput("", "", "")
	if nil != err {
		return err
	}

	l.log.SetOutput(out)
	l.log.AddHook(l.hook)
	l.log.SetLevel(level)
	return nil
}

func (l *Log) initRemoteLogSystem(opts ...option) error {
	ctx := findContext(opts...)
	caller = findCaller(opts...)
	level := findLevel(opts...)

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

func (l *Log) watch(opts ...option) {
	num := findWatchLogsByNum(opts...)
	size := findWatchLogsBySize(opts...)

	var name string = ""
	var tick = time.Tick(time.Millisecond * 500)

	for {
		select {
		case <-tick:
			files, err := ioutil.ReadDir(l.path)
			if err != nil {
				continue
			}

			l.delLogFileByNum(num, files)
			l.cutLogFileBySize(num, size, files)

			if name = filepath.Base(l.logFileName()); l.name != name {
				l.name = name
				atomic.StoreInt32(&l.index, 0)
				l.hook.SetDefaultPath(fmt.Sprintf("%s/%s", l.path, l.name))
			}
		case <-l.stop:
			return
		}
	}

	return
}

func (l *Log) cutLogFileBySize(num int, basic int64, files []os.FileInfo) {
	for _, f := range files {
		if f.IsDir() || l.name != f.Name() {
			continue
		}

		if f.Size() > (basic / int64(num)) {
			name := fmt.Sprintf("%s/%s", l.path, l.name)
			os.Rename(name, fmt.Sprintf("%s_%02d", name, atomic.AddInt32(&l.index, 1)))
			break
		}
	}

	return
}

func (l *Log) delLogFileByNum(num int, files []os.FileInfo) {
	var logs = make(map[string]string)
	var name string
	var timestamps []string = nil

	for _, f := range files {
		if name = f.Name(); !strings.HasPrefix(f.Name(), fmt.Sprintf("./%s", l.path)) {
			name = fmt.Sprintf("./%s/%s", l.path, name)
		}

		if f.IsDir() || !strings.HasPrefix(name, l.formt) {
			continue
		}

		timestamps = append(timestamps, f.ModTime().String())
		logs[f.ModTime().String()] = path.Join(l.path, f.Name())
	}

	sort.Strings(timestamps)

	for i := 0; i < len(timestamps)-num; i++ {
		os.Remove(logs[timestamps[i]])
		l.log.Debugf("remove file:%s.", logs[timestamps[i]])
	}

	return
}
