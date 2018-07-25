package log

import (
	"runtime"

	"github.com/sirupsen/logrus"
)

type entry struct {
	caller bool
	log    *logrus.Entry
}

func (e entry) Debug(args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok || !e.caller {
		e.log.Debug(args...)
		return
	}

	e.log.WithFields(logrus.Fields{"file": file, "line": line}).Debug(args...)
}

func (e entry) Debugf(format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok || !e.caller {
		e.log.Debugf(format, args...)
		return
	}
	e.log.WithFields(logrus.Fields{"file": file, "line": line}).Debugf(format, args...)
}

func (e entry) Info(args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok || !e.caller {
		e.log.Info(args...)
		return
	}
	e.log.WithFields(logrus.Fields{"file": file, "line": line}).Info(args...)
}

func (e entry) Infof(format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok || !e.caller {
		e.log.Infof(format, args...)
		return
	}
	e.log.WithFields(logrus.Fields{"file": file, "line": line}).Infof(format, args...)
}

func (e entry) Error(args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok || !e.caller {
		e.log.Error(args...)
		return
	}
	e.log.WithFields(logrus.Fields{"file": file, "line": line}).Error(args...)
}

func (e entry) Errorf(format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok || !e.caller {
		e.log.Errorf(format, args...)
		return
	}
	e.log.WithFields(logrus.Fields{"file": file, "line": line}).Errorf(format, args...)
}

func (e entry) Warn(args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok || !e.caller {
		e.log.Warn(args...)
		return
	}
	e.log.WithFields(logrus.Fields{"file": file, "line": line}).Warn(args...)
}

func (e entry) Warnf(format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok || !e.caller {
		e.log.Warnf(format, args...)
		return
	}
	e.log.WithFields(logrus.Fields{"file": file, "line": line}).Warnf(format, args...)
}

func (e entry) Fatal(args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok || !e.caller {
		e.log.Fatal(args...)
		return
	}
	e.log.WithFields(logrus.Fields{"file": file, "line": line}).Fatal(args...)
}

func (e entry) Fatalf(format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok || !e.caller {
		e.log.Fatalf(format, args...)
		return
	}

	e.log.WithFields(logrus.Fields{"file": file, "line": line}).Fatalf(format, args...)
}

func (e entry) Panic(args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok || !e.caller {
		e.log.Panic(args...)
		return
	}

	e.log.WithFields(logrus.Fields{"file": file, "line": line}).Panic(args...)
}

func (e entry) Panicf(format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok || !e.caller {
		e.log.Panicf(format, args...)
		return
	}

	e.log.WithFields(logrus.Fields{"file": file, "line": line}).Panicf(format, args...)
}

func (e entry) WithField(key string, value interface{}) *entry {
	_, file, line, ok := runtime.Caller(1)
	if !ok || !e.caller {
		return &entry{
			log:    e.log.WithField(key, value),
			caller: false,
		}
	}

	return &entry{
		log:    e.log.WithFields(logrus.Fields{"file": file, "line": line, key: value}),
		caller: true,
	}
}

func (e entry) WithFields(fields map[string]interface{}) *entry {
	_, file, line, ok := runtime.Caller(1)
	if !ok || !e.caller {
		return &entry{
			log:    e.log.WithFields(logrus.Fields(fields)),
			caller: false,
		}
	}

	fields["file"] = file
	fields["line"] = line

	return &entry{
		log:    e.log.WithFields(logrus.Fields(fields)),
		caller: true,
	}
}
