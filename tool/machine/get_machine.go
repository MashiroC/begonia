package machine

import (
	"log"
	"runtime"
	"strconv"
)

type GetMachineFunc func(map[string]string) error

type machineInfo struct {
	fs []GetMachineFunc
}

var M = newMachineInfo()

func newMachineInfo() machineInfo {
	var fs []GetMachineFunc
	m := machineInfo{fs: fs}
	m.AddFunc(func(m map[string]string) error {
		m["cpu"] = strconv.Itoa(runtime.GOMAXPROCS(0))
		return nil
	})
	return m
}

func (m *machineInfo) AddFunc(fun GetMachineFunc) {
	m.fs = append(m.fs, fun)

}

func (m *machineInfo) GetMachineInfo() (map[string]string, error) {
	info := make(map[string]string)
	var err error
	for _, f := range m.fs {
		err = f(info)
		if err != nil {
			log.Println(err)
			break
		}
	}
	return info, err
}
