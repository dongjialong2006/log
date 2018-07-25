package logger

import (
	"log/http"
	"log/types"
)

type Logger struct {
	Url     string
	Protocol string
}

func (l *Logger) Write(p []byte) (n int, err error) {
	if "" == l.Url {
		return 0, nil
	}

	if types.HTTP == l.Protocol {
		_, err = http.Post(l.Url, string(p))
	}

	if types.HTTPS == l.Protocol {
		_, err = http.PostWithTLS(l.Url, string(p))
	}

	return
}
