package centerlog

import (
	"fmt"
	"testing"
)

/*
 @Author: as
 @Date: Creat in 16:42 2021/3/25
 @Description: begonia
*/

func TestName(t *testing.T) {

	msg := NewMsg()
	msg.Fields = make(map[string]string)
	msg.ServerName="Echo"
	msg.Fields["server"]="Echo"
	//fmt.Println(msg.PutMsg())
	fmt.Println(msg.GetAllMsg())
	fmt.Println(msg.QueryField(msg.Fields))
}
