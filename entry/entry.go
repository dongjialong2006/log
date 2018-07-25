package entry

import (
	"runtime"

	"github.com/sirupsen/logrus"
)

type Entry struct {
	Caller bool
	Log    *logrus.Entry
}

func (e Entry) Debug(args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok || !e.Caller {
		e.Log.Debug(args...)
		return
	}

	e.Log.WithFields(logrus.Fields{"file": file, "line": line}).Debug(args...)
}

func (e Entry) Debugf(format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok || !e.Caller {
		e.Log.Debugf(format, args...)
		return
	}
	e.Log.WithFields(logrus.Fields{"file": file, "line": line}).Debugf(format, args...)
}

func (e Entry) Info(args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok || !e.Caller {
		e.Log.Info(args...)
		return
	}
	e.Log.WithFields(logrus.Fields{"file": file, "line": line}).Info(args...)
}

func (e Entry) Infof(format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok || !e.Caller {
		e.Log.Infof(format, args...)
		return
	}
	e.Log.WithFields(logrus.Fields{"file": file, "line": line}).Infof(format, args...)
}

func (e Entry) Error(args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok || !e.Caller {
		e.Log.Error(args...)
		return
	}
	e.Log.WithFields(logrus.Fields{"file": file, "line": line}).Error(args...)
}

func (e Entry) Errorf(format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok || !e.Caller {
		e.Log.Errorf(format, args...)
		return
	}
	e.Log.WithFields(logrus.Fields{"file": file, "line": line}).Errorf(format, args...)
}

func (e Entry) Warn(args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok || !e.Caller {
		e.Log.Warn(args...)
		return
	}
	e.Log.WithFields(logrus.Fields{"file": file, "line": line}).Warn(args...)
}

func (e Entry) Warnf(format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok || !e.Caller {
		e.Log.Warnf(format, args...)
		return
	}
	e.Log.WithFields(logrus.Fields{"file": file, "line": line}).Warnf(format, args...)
}

func (e Entry) Fatal(args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok || !e.Caller {
		e.Log.Fatal(args...)
		return
	}
	e.Log.WithFields(logrus.Fields{"file": file, "line": line}).Fatal(args...)
}

func (e Entry) Fatalf(format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok || !e.Caller {
		e.Log.Fatalf(format, args...)
		return
	}

	e.Log.WithFields(logrus.Fields{"file": file, "line": line}).Fatalf(format, args...)
}

func (e Entry) Panic(args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok || !e.Caller {
		e.Log.Panic(args...)
		return
	}

	e.Log.WithFields(logrus.Fields{"file": file, "line": line}).Panic(args...)
}

func (e Entry) Panicf(format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok || !e.Caller {
		e.Log.Panicf(format, args...)
		return
	}

	e.Log.WithFields(logrus.Fields{"file": file, "line": line}).Panicf(format, args...)
}

func (e Entry) WithField(key string, value interface{}) *Entry {
	_, file, line, ok := runtime.Caller(1)
	if !ok || !e.Caller {
		return &Entry{
			Log:    e.Log.WithField(key, value),
			Caller: false,
		}
	}

	return &Entry{
		Log:    e.Log.WithFields(logrus.Fields{"file": file, "line": line, key: value}),
		Caller: true,
	}
}

func (e Entry) WithFields(fields map[string]interface{}) *Entry {
	_, file, line, ok := runtime.Caller(1)
	if !ok || !e.Caller {
		return &Entry{
			Log:    e.Log.WithFields(logrus.Fields(fields)),
			Caller: false,
		}
	}

	fields["file"] = file
	fields["line"] = line

	return &Entry{
		Log:    e.Log.WithFields(logrus.Fields(fields)),
		Caller: true,
	}
}
