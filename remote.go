package log

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/sirupsen/logrus"
)

func handle(ctx context.Context, addr string, protocol string) {
	var err error = nil
	var status bool = true
	var conn net.Conn = nil
	var tick = time.Tick(time.Second)

	if err = setOutput(addr, protocol); nil != err {
		fmt.Println(err)
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-tick:
			if conn, err = dail(addr, protocol); nil != err {
				fmt.Println(err)
				if status {
					status = false
				}
				tick = time.Tick(time.Second * time.Duration(1+rand.Intn(6)))
				continue
			}
			if nil != conn {
				conn.Close()
				conn = nil
			}

			if !status {
				if err = setOutput(addr, protocol); nil != err {
					fmt.Println(err)
				} else {
					status = true
				}
			}
		}
	}
	return
}

func setOutput(addr string, protocol string) error {
	switch protocol {
	case TCP, UDP:
		conn, err := dail(addr, protocol)
		if nil != err {
			return err
		}
		logrus.SetOutput(conn)
	case HTTP, HTTPS:
		logrus.SetOutput(&output{
			url:      addr,
			protocol: protocol,
		})
	default:
		return fmt.Errorf("do not persist protocol type:%s.", protocol)
	}

	return nil
}

func dail(addr string, protocol string) (net.Conn, error) {
	conn, err := net.DialTimeout(protocol, addr, time.Second*2)
	if nil != err {
		return nil, err
	}

	return conn, nil
}
