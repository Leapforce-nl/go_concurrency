package concurrency

import (
	"io/ioutil"
	"os"
	"strings"

	errortools "github.com/leapforce-libraries/go_errortools"
)

const (
	defaultFileName string = "__concurrency__"
	separator       string = ";"
)

type Concurrency struct {
	fileName string
}

func NewConcurrency(fileName *string) *Concurrency {
	_fileName := defaultFileName
	if fileName != nil {
		_fileName = *fileName
	}

	return &Concurrency{_fileName}
}

func (c *Concurrency) TryStartProcess(processName string) (bool, *errortools.Error) {
	running, e := c.getRunning()
	if e != nil {
		return false, e
	}

	_running := *running
	isRunning := false

	for _, r := range *running {
		if r == processName {
			isRunning = true
			break
		}
	}

	if !isRunning {
		_running = append(_running, processName)
		e := c.setRunning(_running)
		if e != nil {
			return false, e
		}
	}

	return !isRunning, nil
}

func (c *Concurrency) StopProcess(processName string) *errortools.Error {
	running, e := c.getRunning()
	if e != nil {
		return e
	}

	_running := []string{}
	for _, r := range *running {
		if r != processName {
			_running = append(_running, r)
		}
	}

	e = c.setRunning(_running)
	if e != nil {
		return e
	}

	return nil
}

func (c *Concurrency) getRunning() (*[]string, *errortools.Error) {
	running := []string{}

	fileExists, e := fileExists(c.fileName)
	if e != nil {
		return nil, e
	}

	if fileExists {
		b, err := ioutil.ReadFile(c.fileName)
		if err != nil {
			return nil, errortools.ErrorMessage(err)
		}

		running = strings.Split(string(b), separator)
	}

	return &running, nil
}

func (c *Concurrency) setRunning(running []string) *errortools.Error {
	b := []byte(strings.Join(running, separator))
	err := ioutil.WriteFile(c.fileName, b, 0644)
	if err != nil {
		return errortools.ErrorMessage(err)
	}
	return nil
}

func fileExists(filename string) (bool, *errortools.Error) {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false, nil
	}
	return !info.IsDir(), nil
}
