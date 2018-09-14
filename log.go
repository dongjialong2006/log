package log

import (
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
)

const (
	LOCAL  = "local"
	REMOTE = "remote"
)

var rw sync.RWMutex
var def *Log = nil

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
		log: logrus.New(),
	}

	if "" != findRemoteAddr(opts...) {
		if err := log.initRemoteLog(opts...); nil != err {
			return nil, err
		}
		return log, nil
	}

	if err := log.initLocalLog(name, opts...); nil != err {
		return nil, err
	}

	return log, nil
}
