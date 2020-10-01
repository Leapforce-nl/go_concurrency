package concurrency

import (
	"io/ioutil"
	"os"
	"strings"
)

const (
	defaultFileName string = "concurrency"
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

func (c *Concurrency) TryStartProcess(processName string) (bool, error) {
	running, err := c.getRunning()
	if err != nil {
		return false, err
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
		err := c.setRunning(_running)
		if err != nil {
			return false, err
		}
	}

	return !isRunning, nil
}

func (c *Concurrency) StopProcess(processName string) error {
	running, err := c.getRunning()
	if err != nil {
		return err
	}

	_running := []string{}
	for _, r := range *running {
		if r != processName {
			_running = append(_running, r)
		}
	}

	err = c.setRunning(_running)
	if err != nil {
		return err
	}

	return nil
}

func (c *Concurrency) getRunning() (*[]string, error) {
	running := []string{}

	if fileExists(c.fileName) {
		b, err := ioutil.ReadFile(c.fileName)
		if err != nil {
			return nil, err
		}

		running = strings.Split(string(b), separator)
	}

	return &running, nil
}

func (c *Concurrency) setRunning(running []string) error {
	b := []byte(strings.Join(running, separator))
	err := ioutil.WriteFile(c.fileName, b, 0644)
	return err
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
