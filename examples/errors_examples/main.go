package main

import (
	"fmt"
)

type ConnError struct {
	Line int
	Message string
	Code int
}

func (c ConnError)Error() string {
	return fmt.Sprintf("Message: %s, Line: %d, Code: %d", c.Message, c.Line, c.Code)
}

func connError() error {
	return &ConnError{
		Line: 32,
		Message: "some error",
		Code: 500,
	}
}

func main(){
	if err := connError(); err != nil {
        err, ok := err.(*ConnError)
		if ok {
			fmt.Println(err.Error())
		} else {
			panic(err)
		}
	}
}
