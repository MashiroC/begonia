package main

const (
	workLimit = 50
	nodeNums  = 5
	workNums  = 1000000
)

func main(){
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
					SayHello("shiina")
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
}