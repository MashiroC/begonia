package centerlog

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {

	msg := NewMsg()
	msg.Fields = make(map[string]string)
	msg.ServerName="Echo"
	msg.Fields["server"]="Echo"
	//fmt.Println(msg.PutMsg())
	fmt.Println(msg.GetAllMsg())
	fmt.Println(msg.QueryField(msg.Fields))
}
