package frame

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"
)

func TestPing_Marshal(t *testing.T) {
	p := NewPing(5 * time.Second)
	b := p.Marshal()
	ping, err := unMarshalPing(b)
	fmt.Println(b, ping, err)
}
func TestNewPing(t *testing.T) {
	timer := time.NewTimer(1 * time.Second)
	ctx, cancelFunc := context.WithCancel(context.Background())
	go whatTime(timer, ctx)
	timer = time.NewTimer(5 * time.Second)
	cancelFunc()
	time.Sleep(2 * time.Second)
}
func whatTime(timer *time.Timer, ctx context.Context) {
	<-ctx.Done()
	if timer == nil {
		log.Println("err")
	}
	<-timer.C
	fmt.Println("ok")
}
