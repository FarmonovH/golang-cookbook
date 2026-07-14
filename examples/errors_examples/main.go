package main

import (
	"errors"
	"fmt"
)

type CustomError struct {}

func (c CustomError) Error() string {
	return fmt.Sprintf("some error %d", 404)
}

func unwrap(err error) error {
	for {
		newErr := errors.Unwrap(err)
		if newErr == nil {
			return err
		}
		err = newErr
	}
} 


func main(){
	err1 := fmt.Errorf("some error")
	err3 := fmt.Errorf("err3 %w", err1)
	err4 := fmt.Errorf("err4 %w", err3)
	err5 := fmt.Errorf("err5 %w", err4)
	
	// err2 := errors.New("some error")
	// fmt.Println(err3.Error() == err2.Error())
	fmt.Println(errors.Is(err3, err1))
	// fmt.Println(err3.Error())
	// fmt.Println(err2.Error())

	// fmt.Println(err3)
	// fmt.Println(errors.Unwrap(err3))
	// fmt.Println(errors.Unwrap(err1) == nil)
	val1 := unwrap(err5)
	val2 := err1
	fmt.Println(val1 == val2)


	// err32 := fmt.Errorf("some error")


}
