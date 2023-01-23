package main

import (
	"fmt"
)

type A struct {
	val int
}

func (a *A) update(b int) *A {
	a.val = b
	return a
}

func main() {
	fmt.Println(*((&A{}).update(3)))

}
