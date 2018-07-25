package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const filePathEmpty = "file path is empty."

func ReadFile(name string) ([]byte, error) {
	if "" == name {
		return nil, fmt.Errorf(filePathEmpty)
	}

	_, err := os.Stat(name)
	if nil != err {
		return nil, err
	}

	return ioutil.ReadFile(name)
}

func WriteFile(path string, data []byte) error {
	if "" == path {
		return fmt.Errorf("path is empty.")
	}

	if 0 == len(data) {
		return fmt.Errorf("data len is zero.")
	}

	var dir string = ""
	pos := strings.LastIndex(path, "/")
	if -1 == pos {
		dir = "./" + path
	} else {
		dir = path[:pos+1]
	}

	_, err := os.Stat(dir)
	if nil != err {
		if os.IsNotExist(err) {
			err = os.MkdirAll(dir, os.ModePerm)
			if nil != err {
				return err
			}
		} else {
			return err
		}
	}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if nil != err {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	return err
}

func CreatePath(path string) error {
	if "" == path {
		return fmt.Errorf(filePathEmpty)
	}

	pos := strings.LastIndex(path, "/")
	if -1 == pos {
		return nil
	}
	path = path[:pos]
	_, err := os.Stat(path)
	if nil == err {
		return nil
	}

	if os.IsNotExist(err) {
		err = os.MkdirAll(path, os.ModePerm)
	}

	return err
}

func CheckFileExist(path string) (bool, error) {
	if "" == path {
		return false, fmt.Errorf(filePathEmpty)
	}

	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func DeletePath(path string) error {
	if "" == path {
		return fmt.Errorf(filePathEmpty)
	}

	tmp := strings.Split(path, "/")
	for i := len(tmp); i > 1; i-- {
		os.RemoveAll(strings.Join(tmp[:i], "/"))
	}

	return nil
}
