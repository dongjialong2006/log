package http

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
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

func response(code int, value string) (string, error) {
	if http.StatusOK != code && http.StatusCreated != code {
		if "" != value {
			return value, fmt.Errorf("http response code:%v, error:%s.", code, value)
		}
		return "", fmt.Errorf("http response error code:%v.", code)
	}

	return value, nil
}

func Post(url string, value string) (string, error) {
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

func PostWithTLS(url string, value string) (string, error) {
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
