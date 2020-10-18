package main

import (
	"begonia2/dispatch/conn"
	"fmt"
)

func main() {

	center.New()
	accept, err := conn.Listen(":12306")

	fmt.Println("start")
out:
	for {
		select {
		case c := <-accept:
			fmt.Println("new conn")
			go func(c conn.Conn) {
				for {
					opcode, data, err := c.Recv()
					if err != nil {
						fmt.Println(err)
						return
					}
					fmt.Println(opcode)
					fmt.Println(data)
					fmt.Println()
				}
			}(c)
		case err := <-err:
			fmt.Println(err)
			break out

		}
	}

}
