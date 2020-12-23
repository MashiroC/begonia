package frame

import (
	"fmt"
	"github.com/MashiroC/begonia/tool/machine"
	"testing"
	"time"
)

func TestPing_Marshal(t *testing.T) {
	p := NewPing(5 * time.Second)
	b := p.Marshal()
	ping, err := unMarshalPing(b)
	fmt.Println(b, ping, err)
}
func TestNewPong(t *testing.T) {

	info, err := machine.M.GetMachineInfo()
	machine.M.GetMachineInfo()
	fmt.Println(info, err)
}
