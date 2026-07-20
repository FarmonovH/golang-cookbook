package main

import (
	"fmt"
	"testing"
)

func TestIs(t *testing.T){
	err1 := fmt.Errorf("some error")
	err3 := fmt.Errorf("err3 %w", err1)
	err4 := fmt.Errorf("err4 %w", err3)
	err5 := fmt.Errorf("err5 %w", err4)
	if !Is(err5, err1) {
		t.Fail()
	}
}