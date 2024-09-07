package main

import (
	"02/merdge"
	"fmt"
)

func main() {

	ch1 := make(chan interface{})
	ch2 := make(chan interface{})
	ch3 := make(chan interface{})

	go func() {
		defer close(ch1)

		for i := 0; i < 10; i++ {
			ch1 <- i
		}
	}()
	go func() {
		defer close(ch2)
		for _, c := range "abcdef" {
			ch2 <- string(c)
		}
	}()
	go func() {
		defer close(ch3)
		for i := 20; i < 30; i++ {
			ch3 <- i
		}
	}()

	res := merdge.Multiplex(ch3, ch2, ch1)

	for c := range res {
		fmt.Println(c)
	}


}
