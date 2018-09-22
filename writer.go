package log

import (
	"os"
	"sync"
)

type output struct {
	sync.RWMutex
	URL  string
	Type string
	File *os.File
}

func newOutput(name string, url, stype string) (*output, error) {
	if "" == name {
		return &output{
			URL:  url,
			Type: stype,
			File: nil,
		}, nil
	}
	file, err := os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if nil != err {
		return nil, err
	}

	return &output{
		URL:  url,
		Type: stype,
		File: file,
	}, nil
}

func (l *output) Close() {
	l.Lock()
	if nil != l.File {
		l.File.Close()
	}
	l.Unlock()
}

func (l *output) Write(p []byte) (n int, err error) {
	l.RLock()
	defer l.RUnlock()
	if nil != l.File {
		n, err = l.File.Write(p)
		return
	}

	if HTTP == l.Type {
		_, err = post(l.URL, string(p))
		return
	}

	if HTTPS == l.Type {
		_, err = postWithTLS(l.URL, string(p))
	}

	return
}
