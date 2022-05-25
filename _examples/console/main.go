package main

import (
	"errors"
	"fmt"

	"github.com/mechiru/stackerr"
)

func main() {
	simple()
	wrap()
}

func simple() {
	err := stackerr.New("error")

	fmt.Println("------ error without stack ------")
	fmt.Println(err.Error())
	fmt.Println("------ error with stack ------")
	fmt.Printf("%+v\n", err)
}

func wrap() {
	e1 := errors.New("error1")
	e2 := stackerr.Errorf("error2: %w", e1)

	fmt.Println("------ error without stack ------")
	fmt.Println(e2.Error())
	fmt.Println("------ error with stack ------")
	fmt.Printf("%+v\n", e2)

	fmt.Printf("errors.Is: %v\n", errors.Is(e2, e1))
}
