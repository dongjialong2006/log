package log

import (
	"runtime"

	"github.com/sirupsen/logrus"
)

type Fields map[string]interface{}

type Entry struct {
	caller bool
	log    *logrus.Entry
}

func (e Entry) Debug(args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok || !e.caller {
		e.log.Debug(args...)
		return
	}

	e.log.WithFields(logrus.Fields{"file": file, "line": line}).Debug(args...)
}

func (e Entry) Debugf(format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok || !e.caller {
		e.log.Debugf(format, args...)
		return
	}
	e.log.WithFields(logrus.Fields{"file": file, "line": line}).Debugf(format, args...)
}

func (e Entry) Info(args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok || !e.caller {
		e.log.Info(args...)
		return
	}
	e.log.WithFields(logrus.Fields{"file": file, "line": line}).Info(args...)
}

func (e Entry) Infof(format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok || !e.caller {
		e.log.Infof(format, args...)
		return
	}
	e.log.WithFields(logrus.Fields{"file": file, "line": line}).Infof(format, args...)
}

func (e Entry) Error(args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok || !e.caller {
		e.log.Error(args...)
		return
	}
	e.log.WithFields(logrus.Fields{"file": file, "line": line}).Error(args...)
}

func (e Entry) Errorf(format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok || !e.caller {
		e.log.Errorf(format, args...)
		return
	}
	e.log.WithFields(logrus.Fields{"file": file, "line": line}).Errorf(format, args...)
}

func (e Entry) Warn(args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok || !e.caller {
		e.log.Warn(args...)
		return
	}
	e.log.WithFields(logrus.Fields{"file": file, "line": line}).Warn(args...)
}

func (e Entry) Warnf(format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok || !e.caller {
		e.log.Warnf(format, args...)
		return
	}
	e.log.WithFields(logrus.Fields{"file": file, "line": line}).Warnf(format, args...)
}

func (e Entry) Fatal(args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok || !e.caller {
		e.log.Fatal(args...)
		return
	}
	e.log.WithFields(logrus.Fields{"file": file, "line": line}).Fatal(args...)
}

func (e Entry) Fatalf(format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok || !e.caller {
		e.log.Fatalf(format, args...)
		return
	}

	e.log.WithFields(logrus.Fields{"file": file, "line": line}).Fatalf(format, args...)
}

func (e Entry) Panic(args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok || !e.caller {
		e.log.Panic(args...)
		return
	}

	e.log.WithFields(logrus.Fields{"file": file, "line": line}).Panic(args...)
}

func (e Entry) Panicf(format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok || !e.caller {
		e.log.Panicf(format, args...)
		return
	}

	e.log.WithFields(logrus.Fields{"file": file, "line": line}).Panicf(format, args...)
}

func (e Entry) WithField(key string, value interface{}) *Entry {
	return &Entry{
		log:    e.log.WithField(key, value),
		caller: e.caller,
	}
}

func (e Entry) WithFields(fields Fields) *Entry {
	return &Entry{
		log:    e.log.WithFields(logrus.Fields(fields)),
		caller: e.caller,
	}
}

func (e Entry) WithError(err error) *Entry {
	return &Entry{
		log:    e.log.WithField("error", err),
		caller: e.caller,
	}
}
