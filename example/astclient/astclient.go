package main

import (
	"fmt"
	"github.com/MashiroC/begonia/example/server/call"
)

const (
	workLimit = 50
	nodeNums  = 5
	workNums  = 1000000
)

func main() {
	res, err := call.SayHello("kieran")
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

//fmt.Println(SayHello("shiina"))
//fmt.Println(Add(1,2))
//fmt.Println(Mod(5,2))
//wg1 := sync.WaitGroup{}
//
//for i := 0; i < nodeNums; i++ {
//	wg1.Add(1)
//	go func() {
//		ch := make(chan struct{}, workLimit)
//		for i := 0; i < workLimit; i++ {
//			ch <- struct{}{}
//		}
//
//		t := time.Now()
//
//		wg := sync.WaitGroup{}
//		for i := 0; i < workNums; i++ {
//			<-ch
//			wg.Add(1)
//			go func() {
//				defer func() {
//					wg.Done()
//					ch <- struct{}{}
//				}()
//i1, i2, i3, i4, i5, f1, f2, ok, str, s1, s2, s6, st, m1, m2, m3 := Echo2()
//fmt.Println(i1, i2, i3, i4, i5, f1, f2, ok)
//fmt.Println(str)
//fmt.Println(s1, s2, s6)
//fmt.Println(st, m1, m2, m3)
//			}()
//		}
//
//		wg.Wait()
//
//		fmt.Println(time.Now().Sub(t).String())
//
//		wg1.Done()
//	}()
//}
//wg1.Wait()
//}
