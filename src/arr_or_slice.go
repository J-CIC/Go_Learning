package main

import (
	"fmt"
	"reflect"
)

func main() {
	var type1 = [5]int{1, 2, 3, 4, 5}
	var type2 = [...]int{1, 2, 3, 4, 5}
	var type3 = []int{1, 2, 3, 4, 5}
	type4 := []int{1, 2, 3, 4, 5}
	var type5 = make([]int, 3, 5)
	var type6 = new([5]int)[0:3]

	fmt.Printf("%T,%s,%v\n", type1, reflect.TypeOf(type1).Kind(), type1)
	// [5]int,array,[1 2 3 4 5]

	fmt.Printf("%T,%s,%v\n", type2, reflect.TypeOf(type2).Kind(), type2)
	// [5]int,array,[1 2 3 4 5]

	fmt.Printf("%T,%s,%v\n", type3, reflect.TypeOf(type3).Kind(), type3)
	// []int,slice,[1 2 3 4 5]

	fmt.Printf("%T,%s,%v\n", type4, reflect.TypeOf(type4).Kind(), type4)
	// []int,slice,[1 2 3 4 5]

	fmt.Printf("%T,%s,%v\n", type5, reflect.TypeOf(type5).Kind(), type5)
	// []int,slice,[0 0 0]

	fmt.Printf("%T,%s,%v\n", type6, reflect.TypeOf(type6).Kind(), type6)
	// []int,slice,[0 0 0]
}
