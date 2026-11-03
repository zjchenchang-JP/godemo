package main

import (
	"fmt"
	"time"
)

func printCh(u <-chan int) {
	time.Sleep(2 * time.Second)
	fmt.Println("print int", <-u)
}

func main() {
	// c := make(chan int, 5)
	// a := 0
	// c <- a
	// fmt.Println(a)
	// // modify
	// a = 1

	// go printCh(c)
	// time.Sleep(5 * time.Second)
	// fmt.Println(a)

	// 实现限流器
	// 每次处理3个请求
	chLimit := make(chan struct{}, 5)
	for i := 0; i < 20; i++ {
		chLimit <- struct{}{}
		go func(i int) {
			fmt.Println("下游处理服务器...", i)
			time.Sleep(time.Second * 2)
			<-chLimit
		}(i)
	}
	time.Sleep(20 * time.Second)
}

// 用channel实现一个互斥锁

