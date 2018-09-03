package log

import (
	"github.com/sirupsen/logrus"
)

const (
	LOCAL  = "local"
	REMOTE = "remote"
)

func New(name string) *Entry {
	return logger.WithField("model", name)
}

func NewLog(name string, opts ...option) (*Log, error) {
	log := &Log{
		log: logrus.New(),
	}

	if err := log.initLocalLog(name, opts...); nil != err {
		return nil, err
	}

	return log, nil
}
