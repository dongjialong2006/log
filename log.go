package log

import (
	"fmt"
	"sync"

	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

const (
	LOCAL  = "local"
	REMOTE = "remote"
)

var rw sync.RWMutex

var def *Log = nil
var opt []option = nil

func init() {
	opt = append(opt, WithLogLevel("debug"))
	opt = append(opt, WithLogName("./log/system.log"))
	opt = append(opt, WithWatchEnable(false))
	opt = append(opt, WithTerminal(false))
}

type Log struct {
	log   *logrus.Logger
	self  bool
	path  string
	name  string
	stop  chan struct{}
	hook  *lfshook.LfsHook
	index int32
	formt string
}

func (l *Log) NewEntry(name string) *Entry {
	var e *Entry

	if "" == name {
		e = &Entry{
			log:    l.log,
			caller: caller,
		}
	} else {
		e = &Entry{
			log:    l.log.WithField("model", name),
			caller: caller,
		}
	}

	for key, field := range fields {
		e.log = e.log.WithField(key, field)
	}

	return e
}

func (l *Log) Stop() {
	if nil == l.stop {
		return
	}

	select {
	case <-l.stop:
	default:
		close(l.stop)
	}
}

func New(name string, opts ...option) *Entry {
	rw.Lock()
	defer rw.Unlock()

	if nil == def {
		var err error = nil
		def, err = NewLog("", opts...)
		if nil != err {
			fmt.Println(err)
			return nil
		}
	}

	return def.NewEntry(name)
}

func NewLog(name string, opts ...option) (*Log, error) {
	log := &Log{
		log:  logrus.New(),
		stop: make(chan struct{}),
	}

	if 0 == len(opts) {
		opts = opt
	}

	if "" != findRemoteAddr(opts...) {
		if err := log.initRemoteLogSystem(opts...); nil != err {
			return nil, err
		}
		return log, nil
	}

	if err := log.initLocalLogSystem(name, opts...); nil != err {
		return nil, err
	}

	return log, nil
}
