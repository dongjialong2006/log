package log

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

var defaultClient = &http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		MaxIdleConnsPerHost: 1,
	},
	Timeout: 30 * time.Second,
}

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
		log.Infof("new connect was created from %s to %s.", conn.LocalAddr(), conn.RemoteAddr())
		logrus.SetOutput(conn)
	case HTTP, HTTPS:
		logrus.SetOutput(&logger{
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

func response(code int, value string) (string, error) {
	if http.StatusOK != code && http.StatusCreated != code {
		if "" != value {
			return value, fmt.Errorf("http response code:%v, error:%s.", code, value)
		}
		return "", fmt.Errorf("http response error code:%v.", code)
	}

	return value, nil
}

func post(url string, value string) (string, error) {
	if "" == url {
		return "", fmt.Errorf("url is empty.")
	}

	var req *http.Request = nil
	var err error = nil
	if "" == value {
		req, err = http.NewRequest("POST", url, nil)
		if nil != err {
			return "", err
		}
	} else {
		req, err = http.NewRequest("POST", url, strings.NewReader(value))
		if nil != err {
			return "", err
		}
	}

	resp, err := defaultClient.Do(req)
	if nil != err {
		return "", err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		return "", err
	}

	return response(resp.StatusCode, string(data))
}

func postWithTLS(url string, value string) (string, error) {
	if "" == url {
		return "", fmt.Errorf("url is empty.")
	}
	resp, err := defaultClient.Post(url, "application/json", strings.NewReader(value))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		return "", err
	}

	return response(resp.StatusCode, string(data))
}
