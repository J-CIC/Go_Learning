package main

import (
	"fmt"
)

func main() {
	slice := make([]int, 0, 4)
	for i, _ := range slice {
		slice = append(slice,i)
	}
	fmt.Printf("%+v\n", slice)
	n := len(slice)
	a := slice[n:n]
	fmt.Printf("len of [n:n]: %d,%v\n",len(a),a)
	a = slice[n:n+1]
	fmt.Printf("len of [n:n+1]: %d,%v\n",len(a),a)
}
