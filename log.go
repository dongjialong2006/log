package log

import (
	"fmt"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

const (
	LOCAL  = "local"
	REMOTE = "remote"
)

var rw sync.RWMutex
var def *Log = nil

type Log struct {
	rw     sync.RWMutex
	log    *logrus.Logger
	self   bool
	name   string
	file   *os.File
	index  int32
	source string
}

func (r *Log) NewEntry(name string) *Entry {
	e := &Entry{
		log:    r.log.WithField("model", name),
		caller: caller,
	}

	for key, field := range fields {
		e.log = e.log.WithField(key, field)
	}

	return e
}

func New(name string, opts ...option) *Entry {
	rw.Lock()
	defer rw.Unlock()
	if nil == def {
		if 0 == len(opts) {
			opts = append(opts, WithLogLevel("debug"))
			opts = append(opts, WithLogName("./log/system.log"))
			opts = append(opts, WithWatchEnable(false))
			opts = append(opts, WithTerminal(false))
		}

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
		file: nil,
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
