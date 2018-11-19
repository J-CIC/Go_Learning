package main

import (
	"fmt"
	"runtime"
)

func main() {
	println("Hello world")
	fmt.Printf("%s", runtime.Version())
}
