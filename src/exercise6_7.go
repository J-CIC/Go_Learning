package main

import (
	"fmt"
	"strings"
)

func main() {
	str := "this is 汉字"
	str = strings.Map(func(r rune) rune{
		if r > 255 {
			return '?'
		}
		return r
	},str)
	fmt.Printf(str)
}
