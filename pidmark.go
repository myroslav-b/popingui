package main

import (
	"errors"
	"io/ioutil"
	"os"
	"strconv"
)

const (
	cPidFileName = "pid.popingui"
)

type tPidMark struct {
	dirName  string
	fileName string
	pid      int
}

func initPidMark(dirName, fileName string, pid int) (tPidMark, error) {
	var pidMark tPidMark
	pidMark.dirName = dirName
	pidMark.fileName = fileName
	pidMark.pid = 0
	data, err := ioutil.ReadFile(dirName + "/" + fileName)
	if err == nil {
		oldPid, err := strconv.Atoi(string(data))
		if err != nil {
			pidMark.pid = -1
			return pidMark, errors.New("The PID file already exist, but is contents are incorrect")
		}
		pidMark.pid = oldPid
		return pidMark, nil
	}
	return pidMark, nil
}

func (pidMark tPidMark) new(pid int) error {
	pidFile, err := os.Create(pidMark.dirName + "/" + pidMark.fileName)
	if err != nil {
		return errors.New("Can't create PID file")
	}
	defer pidFile.Close()
	_, err = pidFile.WriteString(strconv.Itoa(pid))
	if err != nil {
		return errors.New("Can't write PID file")
	}
	pidMark.pid = pid
	return nil
}

func (pidMark tPidMark) getPid() int {
	return pidMark.pid
}

func (pidMark tPidMark) remove() error {
	defer func() {
		pidMark.dirName = ""
		pidMark.fileName = ""
		pidMark.pid = -2
	}()
	err := os.Remove(pidMark.dirName + "/" + pidMark.fileName)
	if err != nil {
		return errors.New("Can't delete PID file")
	}
	return nil
}
