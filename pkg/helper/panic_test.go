package helper

import (
	"errors"
	"fmt"
	"testing"
)

func restore() {
	if r := recover(); r != nil {
		fmt.Printf("Recovered in f: %v\n", r)
	}
}

// go test -run TestErrPanicIfErr
func TestErrPanicIfErr(t *testing.T) {
	defer restore()

	for i := 0; i < 99999; i++ {
		PanicIfErr(errors.New("api error"))
		t.Errorf("The code did not panic")
	}
}

// go test -run TestNilPanicIfErr
func TestNilPanicIfErr(t *testing.T) {
	defer restore()

	for i := 0; i < 99999; i++ {
		PanicIfErr(nil)
	}
}
