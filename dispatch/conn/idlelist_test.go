package conn

import (
	"fmt"
	"testing"
	"time"
)

func TestIdle(t *testing.T) {
	list := &idleList{
		len:   0,
		front: nil,
		back:  nil,
	}
	for i := 0; i < 5; i++ {
		list.pushFront(&poolConn{
			t: time.Now(),
		})
	}

	fmt.Println(list.len)
	fmt.Println(list.back)
	fmt.Println(list.front)
}
