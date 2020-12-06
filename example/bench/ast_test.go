package bench

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

//func BenchmarkAstEncode(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		AstEncode()
//	}
//}
////
//func BenchmarkAstDecode(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		AstDecode()
//	}
//}
////
//func BenchmarkHambaEncode(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		hambaEncode()
//	}
//}
////
//func BenchmarkHambaDecode(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		hambaDecode()
//	}
//}

//func TestAst(t *testing.T) {
//	hamba()
//	ast()
//}

//func BenchmarkHamba(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		hamba()
//	}
//}

//func BenchmarkAst(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		ast()
//	}
//}

func TestMutex(t *testing.T) {
	//err := Timeout(func() {
	//	time.Sleep(time.Second*100)
	//})
	//if err != nil {
	//	panic(err)
	//}
	var i int32
	flag := atomic.CompareAndSwapInt32(&i, 10, 11)
	fmt.Println(flag)
}

func BenchmarkChan(b *testing.B) {
	//ctx,_:=context.WithTimeout(context.Background(),5*time.Second)
	for i := 0; i < b.N; i++ {
		ch := make(chan struct{})
		go func() {
			ch <- struct{}{}
		}()
		select {
		case <-time.After(10 * time.Second):
			panic("timeout")
		case <-ch:
			// any code
		}
	}
}

func BenchmarkMutex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := Timeout(func() {
			// do some thing
		})
		if err != nil {
			panic(err)
		}
	}
}

func Timeout(f func()) error {
	var pos bool
	var isTimeout bool
	var l, l2 sync.Mutex
	l.Lock()

	time.AfterFunc(3*time.Second, func() {
		l2.Lock()
		if !pos {
			pos = true
			isTimeout = true
			l.Unlock()
		}
		l2.Unlock()
	})

	go func() {
		f()
		l2.Lock()
		if !pos {
			pos = true
			l.Unlock()
		}
		l2.Unlock()
	}()
	l.Lock()
	l.Unlock()
	if isTimeout {
		return errors.New("timeout")
	}
	return nil
}
