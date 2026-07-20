package methods

import (
	"fmt"
	"math/rand"
)

type Bank struct {

}

func NewBank() *Bank {
	return &Bank{}
}

func (b Bank) Pay(description string, price float64) int {
	fmt.Println("payment menthod via with Bank")
	fmt.Printf("Description: %s, Price: %f", description, price)
	return rand.Int()
}


func (b Bank) Cancel(id int) {
	fmt.Printf("canceled %d\n", id)
}