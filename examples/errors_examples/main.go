package main

import (
	"fmt"
)

func  A(){
	defer fmt.Println("defered A")
	B()
}

func B(){
	defer func(){
		if err := recover(); err != nil {
			fmt.Println("panic ", err)
		} else {
			fmt.Println("panic not found")
		}
	}()
	C()
}

func C() {
	defer fmt.Println("defer C")
	raisePanic()
}

func raisePanic(){
	panic("raise panic function")
}

func main() {
	A()
}