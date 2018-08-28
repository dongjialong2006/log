package log

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/sirupsen/logrus"
)

var logger = New("log")

func handle(ctx context.Context, log *logrus.Logger, addr string, protocol string) {
	var tick = time.Tick(time.Second * 3)

	conn, err := setOutput(log, addr, protocol)
	if nil != err {
		logger.Error(err)
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-tick:
			if nil == conn {
				conn, err = setOutput(log, addr, protocol)
				if nil != err {
					logger.Error(err)
					tick = time.Tick(time.Second * time.Duration(1+rand.Intn(6)))
					continue
				}
			}

			if _, err = conn.Write([]byte("keepalive")); nil != err {
				logger.Error(err)
				tick = time.Tick(time.Second * time.Duration(1+rand.Intn(6)))
			}
		}
	}

	return
}

func setOutput(log *logrus.Logger, addr string, protocol string) (net.Conn, error) {
	switch protocol {
	case TCP, UDP:
		conn, err := dail(addr, protocol)
		if nil != err {
			return nil, err
		}
		if nil == log {
			logrus.SetOutput(conn)
		} else {
			log.Out = conn
		}
		return conn, nil
	case HTTP, HTTPS:
		if nil == log {
			logrus.SetOutput(&output{
				url:      addr,
				protocol: protocol,
			})
		}
	default:
		return nil, fmt.Errorf("do not persist protocol type:%s.", protocol)
	}

	return nil, nil
}

func dail(addr string, protocol string) (net.Conn, error) {
	conn, err := net.DialTimeout(protocol, addr, time.Second*2)
	if nil != err {
		return nil, err
	}

	return conn, nil
}
