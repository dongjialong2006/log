package log

type output struct {
	url      string
	protocol string
}

func (l *output) Write(p []byte) (n int, err error) {
	if "" == l.url {
		return 0, nil
	}

	if HTTP == l.protocol {
		_, err = post(l.url, string(p))
	}

	if HTTPS == l.protocol {
		_, err = postWithTLS(l.url, string(p))
	}

	return
}
