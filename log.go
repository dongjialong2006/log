package log

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

const (
	LOCAL  = "local"
	REMOTE = "remote"
)

var def *Log = nil

func New(name string, opts ...option) *Entry {
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
		log:   logrus.New(),
		paths: make(map[string]int),
	}

	if err := log.initLocalLog(name, opts...); nil != err {
		return nil, err
	}

	return log, nil
}
