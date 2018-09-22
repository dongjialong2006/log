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

	"github.com/sirupsen/logrus"
	formatter "github.com/x-cray/logrus-prefixed-formatter"
)

const filePathEmpty = "file path is empty."

func (l *Log) close() {
	if nil != l.file {
		l.file.Close()
	}
}

func (l *Log) logFileName() (string, error) {
	if !l.self {
		name := l.source
		if err := l.newFileDir(name); nil != err {
			return "", err
		}
		year, month, day := time.Now().Date()
		return fmt.Sprintf("%s_%d.%d_%02d_%02d", name, os.Getpid(), year, month, day), nil
	}

	pos := strings.LastIndex(l.source, "/")
	if -1 == pos {
		return fmt.Sprintf("./%s", l.source), nil
	}

	return l.source, nil
}

func (l *Log) initLogrusLog(level logrus.Level) error {
	l.log.SetLevel(level)
	l.log.Formatter = &formatter.TextFormatter{
		TimestampFormat:  "2006-01-02 15:04:05.00000000",
		ForceColors:      true,
		QuoteEmptyFields: true,
		FullTimestamp:    true,
	}

	return nil
}

func (l *Log) initLocalLogSystem(name string, opts ...option) error {
	level := findLevel(opts...)
	caller = findCaller(opts...)

	if findTerminal(opts...) {
		return terminal(level)
	}

	if "" != name {
		l.self = true
		l.source = name
	} else {
		l.source = findLogName(opts...)
	}

	if err := l.hook(); err != nil {
		return err
	}

	if findWatcherEnable(opts...) {
		go l.watch(opts...)
	}

	return l.initLogrusLog(level)
}

func (l *Log) hook() error {
	name, err := l.logFileName()
	if nil != err {
		return err
	}

	if l.checkFileExist(name) {
		return nil
	}

	file, err := newOutput(name, "", "")
	if nil != err {
		return err
	}

	l.log.SetOutput(file)

	l.rw.Lock()
	if name != l.name {
		l.name = name
	}
	l.close()
	l.file = file
	l.rw.Unlock()

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

	var name string = l.name
	var path string = "./"

	pos := strings.LastIndex(l.name, "/")
	if -1 != pos {
		path = l.name[:pos+1]
		name = l.name[pos+1:]
	}
	defer l.close()

	tick := time.Tick(time.Second)
	for {
		<-tick
		files, err := ioutil.ReadDir(path)
		if err != nil {
			continue
		}

		l.delLogFileByNum(name, path, num, files)
		l.cutLogFileBySize(size, files)

		if err = l.hook(); nil != err {
			l.log.Error(err)
		}
	}

	return
}

func (l *Log) newFileDir(path string) error {
	if "" == path {
		return fmt.Errorf(filePathEmpty)
	}

	pos := strings.LastIndex(path, "/")
	if -1 == pos {
		return nil
	}
	path = path[:pos]
	_, err := os.Stat(path)
	if nil == err {
		return nil
	}

	if os.IsNotExist(err) {
		err = os.MkdirAll(path, os.ModePerm)
	}

	return err
}

func (l *Log) checkFileExist(name string) bool {
	_, err := os.Stat(name)
	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}

	l.log.Error(err)

	return false
}

func (l *Log) cutLogFileBySize(basic int64, files []os.FileInfo) {
	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(l.name, f.Name()) {
			continue
		}

		if f.Size() > basic {
			os.Rename(l.name, fmt.Sprintf("%s_%d", l.name, atomic.AddInt32(&l.index, 1)))
			break
		}
	}

	return
}

func (l *Log) delLogFileByNum(name string, dir string, num int, files []os.FileInfo) {
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

	return
}
