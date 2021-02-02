package machine

import (
	"fmt"
	"testing"
	"time"
)

func TestMachine(t *testing.T) {

	m := NewMachine()
	time.Sleep(2 * time.Second)
	//m.chain.Sign(NewCpuMonitor())
	//m.chain.Sign(NewMemMonitor())
	fmt.Println(m.GetMachineInfo(7))
}
