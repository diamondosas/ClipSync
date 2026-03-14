package main

import (
	"fmt"
	"time"
)
var ch = make(chan int)

func main(){
	start := time.Now()
	go count()

	for v := range ch{
		fmt.Println(v)
	}
	end := time.Since(start)
	fmt.Println(end)
}

func count(){
	for i:=0; i<10000; i++{
		ch <- i
	}
	close(ch)
}