package machine

import (
	"github.com/MashiroC/begonia/tool/chain"
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

type Machine struct {
	chain *chain.Chain
}

func (machine *Machine) GetMachineInfo(code byte) map[string]string {
	info := make(map[string]string)
	req := &chain.Request{
		Code:   code,
		ResFun: func(i interface{}) {
			if m, ok := i.(map[string]string); ok {
				for k, v := range m {
					info[k] = v
				}
			}
		},
	}
	machine.chain.Handle(req)

	return info
}

func NewMachine() *Machine {
	machine := &Machine{
		chain: chain.NewChain(),
	}

	machine.chain.Sign(NewCpuMonitor())

	return machine
}


