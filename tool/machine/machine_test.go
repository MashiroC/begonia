package machine

import (
	"fmt"
	"testing"
)

func TestMachine(t *testing.T) {
	m := NewMachine()
	m.chain.Sign(NewCpuMonitor())
	m.chain.Sign(NewMemMonitor())
	fmt.Println(m.GetMachineInfo(7))
}
