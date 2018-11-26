package main

import (
	"fmt"
)

type Speakable interface {
	speak() int
}

type Playable interface {
	play() string
	Speakable // 接口中内嵌接口
}

type cnn struct {
	count int
	text  string
}

type TV struct {
	cnn       // 结构体中内嵌结构体
	Speakable // 结构体中内嵌接口
}

type TV2 struct {
	Speakable
	*cnn // 结构体中内嵌结构体
}

func (tv *TV) speak() int {
	return tv.count
}

func (tv *TV) play() string {
	return tv.text
}

func (tv *TV2) speak() int {
	return tv.count
}

func (tv *TV2) play() string {
	return tv.text
}

func (c *cnn) speak() int {
	return c.count
}

func main() {
	tv := &TV{cnn{count: 1, text: "try first text"}, &cnn{count: 2, text: "try other text"}}
	var inter_var Playable
	inter_var = tv
	fmt.Println("call speak(),", tv.cnn.speak())
	fmt.Println("call speak(),", tv.Speakable.speak())
	fmt.Println("call speak(),", tv.speak())
	fmt.Println("call speak(),", inter_var.speak())
	fmt.Println("call play(),", inter_var.play())

	fmt.Println("----------------------")

	tv2 := &TV2{&cnn{count: 1, text: "try first text"}, &cnn{count: 2, text: "try other text"}}
	var inter_var2 Playable
	inter_var2 = tv2
	fmt.Println("call speak(),", tv2.cnn.speak())
	fmt.Println("call speak(),", tv2.Speakable.speak())
	fmt.Println("call speak(),", tv2.speak())
	fmt.Println("call speak(),", inter_var2.speak())
	fmt.Println("call play(),", inter_var2.play())
}

// 输出结果，可以看出当同级同名时（比较两组的第3-5输出），先声明的会被使用
// call speak(), 1
// call speak(), 2
// call speak(), 1
// call speak(), 1
// call play(), try first text
// ----------------------
// call speak(), 2
// call speak(), 1
// call speak(), 2
// call speak(), 2
// call play(), try other text
