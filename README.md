Go入门学习
===

内容参考自: [《The Way to Go》](https://go.fdos.me/)

目录
---
<details>
<summary>点击展开目录菜单</summary>

<!-- TOC -->

- [Go及环境相关](#Go及环境相关)
- [与其他语言不同之处](#与其他语言不同之处)
- [数组与切片](#数组与切片)
- [Map相关](#Map相关)
- [包(Package)](#包(Package))
- [结构(struct)和方法(method)](#结构和方法)
- [接口](#接口)
<!-- /TOC -->

</details>


## Go及环境相关

<details>
<summary>How to go get with proxy in windows?</summary>

```bash
#本来是有ss和全局代理软件的，但是想试试不用全局怎么做，以下是自身尝试成功的做法（终端为git bash）：
https_proxy=127.0.0.1:1080 http_proxy=127.0.0.1:1080 go get golang.org/x/tour
```

</details>

<details>
    <summary>Go语言(1.5开始)的自举</summary>

首先什么是编程语言的自举？
以下[回答](https://segmentfault.com/q/1010000000692678)来自segmentfault

> 就是自己的编译器可以自行编译自己的编译器。\
> 实现方法就是这个编译器的作者用这个语言的一些特性来编写编译器并在该编译器中支持这些自己使用到的特性。\
> 首先，第一个编译器肯定是用别的语言写的（不论是C还是Go还是Lisp还是Python），后面的版本才能谈及自举。\
> 至于先有鸡还是先有蛋，我可以举个这样的不太恰当的例子：比如我写了一个可以自举的C编译器叫作mycc，不论是编译器本身的执行效率还是生成的代码的质量都远远好于gcc（本故事纯属虚构），但我用的都是标准的C写的，那么我可以就直接用gcc编译mycc的源码，得到一份可以生成高质量代码但本身执行效率低下的mycc，然后当然如果我再用这个生成的mycc编译mycc的源码得到新的一份mycc，新的这份不光会产生和原来那份同等高质量的代码，而且还能拥有比先前版本更高的执行效率（因为前一份是gcc的编译产物，后一份是mycc的编译产物，而mycc生成的代码质量要远好于gcc的）。故事虽然是虚构的，但是道理差不多就是这么个道理。这也就是为什么如果从源码编译安装新版本的gcc的话，往往会“编译——安装”两到三遍的原因。

</details>

<details>
    <summary>简单运行Go程序</summary>

```bash
go run hello_world.go
```

</details>

<details>
    <summary>Go代码风格格式化</summary>

```bash
gofmt -w *.go
gofmt <foldername>
```
</details>

## 与其他语言不同之处

<details>
    <summary>指针不可运算</summary>
对于经常导致 C 语言内存泄漏继而程序崩溃的指针运算（所谓的指针算法，如：pointer+2，移动指针指向字符串的字节数或数组的某个位置）是不被允许的。Go 语言中的指针保证了内存安全，更像是 Java、C# 和 VB.NET 中的引用。

因此```c = *p++```在 Go 语言的代码中是不合法的。
</details>

<details>
    <summary>命名的返回值</summary>
可以通过在函数签名中声明返回值的名字，从而省略return中的变量，example：

```go
func getX2AndX3_2(input int) (x2 int, x3 int) {
    x2 = 2 * input
    x3 = 3 * input
    // return x2, x3
    return
}
```
</details>

<details>
    <summary>变长参数和Printf</summary>

### 同类型的变长参数
首先看看函数中的语法定义

```go
/**
    FunctionType   = "func" Signature .
    Signature      = Parameters [ Result ] .
    Result         = Parameters | Type .
    Parameters     = "(" [ ParameterList [ "," ] ] ")" .
    ParameterList  = ParameterDecl { "," ParameterDecl } .
    ParameterDecl  = [ IdentifierList ] [ "..." ] Type .
*/
func min(s ...int) int {
    if len(s)==0 {
        return 0
    }
    min := s[0]
    for _, v := range s {
        if v < min {
            min = v
        }
    }
    return min
}
// usage
result := min(1,5,4,2,4)
slice := []int{7,9,3,5,1}
result = min(slice...)
```

### 不同类型的变长参数(以Printf为例)

```go
//一个简单的例子
func typecheck(..,..,values … interface{}) {
    for _, value := range values {
        switch v := value.(type) {
            case int: …
            case float: …
            case string: …
            case bool: …
            default: …
        }
    }
}

// 例如fmt.Printf()
// Printf formats according to a format specifier and writes to standard output.
// It returns the number of bytes written and any write error encountered.
func Printf(format string, a ...interface{}) (n int, err error) {
    return Fprintf(os.Stdout, format, a...)
}

// Fprintf formats according to a format specifier and writes to w.
// It returns the number of bytes written and any write error encountered.
func Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error) {
    p := newPrinter()
    p.doPrintf(format, a)
    n, err = w.Write(p.buf)
    p.free()
    return
}

func (p *pp) doPrintf(format string, a []interface{}) {
    end := len(format)
    argNum := 0         // we process one argument per non-trivial format
    afterIndex := false // previous item in format was an index like [3].
    p.reordered = false

    // some source code that handles the format string is omitted here
    // ......
    // some source code that handles the format string is omitted here

    if !p.reordered && argNum < len(a) {
        p.fmt.clearflags()
        p.buf.WriteString(extraString)
        for i, arg := range a[argNum:] {
            if i > 0 {
                p.buf.WriteString(commaSpaceString)
            }
            if arg == nil {
                p.buf.WriteString(nilAngleString)
            } else {
                p.buf.WriteString(reflect.TypeOf(arg).String())
                p.buf.WriteByte('=')
                p.printArg(arg, 'v')
            }
        }
        p.buf.WriteByte(')')
    }
}
```
</details>

<details>
    <summary>defer推迟执行</summary>
关键字 defer 允许我们推迟到函数返回之前（或任意位置执行 return 语句之后）一刻才执行某个语句或函数（为什么要在返回之后才执行这些语句？因为 return 语句同样可以包含一些操作，而不是单纯地返回某个值）。

关键字 defer 的用法类似于面向对象编程语言 Java 和 C# 的 finally 语句块，它一般用于释放某些已分配的资源。

```go
// open a file  
defer file.Close()

// open a database connection  
defer disconnectFromDB()

// 甚至用来调试函数
package main

import "fmt"

func trace(s string)   { fmt.Println("entering:", s) }
func untrace(s string) { fmt.Println("leaving:", s) }

func a() {
    trace("a")
    defer untrace("a")
    fmt.Println("in a")
}

func b() {
    trace("b")
    defer untrace("b")
    fmt.Println("in b")
    a()
}

func main() {
    b()
}

// 更简洁的版本
package main

import "fmt"

func trace(s string) string {
    fmt.Println("entering:", s)
    return s
}

func un(s string) {
    fmt.Println("leaving:", s)
}

func a() {
    defer un(trace("a"))
    fmt.Println("in a")
}

func b() {
    defer un(trace("b"))
    fmt.Println("in b")
    a()
}

func main() {
    b()
}
```
</details>

## 数组与切片

<details>
    <summary>容易混淆的声明</summary>
看代码和运行结果更直观

```go
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
```
</details>

## Map相关

<details>
    <summary>关于Map的初始化</summary>

```go
//var map1 map[keytype]valuetype
var map1 map[string]int

//值作为切片值，应对一key多value的情况
mp1 := make(map[int][]int)
mp2 := make(map[int]*[]int)
```

> 请永远用make来初始化Map，而不是用new，否则你会获得一个空饮用的指针，相当于声明了一个未初始化的变量并且取得了它的地址

</details>

<details>
    <summary>map中不存在的key的value的初始值</summary>

当Key不存在的时候，返回的是valuetype的空值，判断key是否存在的方式如下：

```go
if _, ok := map1[key1]; ok {
    // 如果存在，ok为true
}
```

删除key的时候直接```delete(map,key)```即可，即便key不存在也不会失败
</details>

<details>
    <summary>Map类型的切片</summary>

代码如下：

```go
package main
import "fmt"

func main() {
    // Version A:
    items := make([]map[int]int, 5)
    for i:= range items {
        items[i] = make(map[int]int, 1)
        items[i][1] = 2
    }
    fmt.Printf("Version A: Value of items: %v\n", items)
    //Version A: Value of items: [map[1:2] map[1:2] map[1:2] map[1:2] map[1:2]]


    // Version B: NOT GOOD!
    items2 := make([]map[int]int, 5)
    for _, item := range items2 {
        item = make(map[int]int, 1) // item is only a copy of the slice element.
        item[1] = 2 // This 'item' will be lost on the next iteration.
    }
    fmt.Printf("Version B: Value of items: %v\n", items2)
    //Version B: Value of items: [map[] map[] map[] map[] map[]]

    // B版本中的item只是一个copy，所以不是一个好的实践，也没有办法真正的初始化到map中

}
```

</details>

<details>
    <summary>Map中的排序</summary>

Map中是不排序的，不论key还是value，若要实现排序有两个思路：

1. 取出其中的所有key到切片中，然后再for-range打印：

```go
// the telephone alphabet:
package main
import (
    "fmt"
    "sort"
)

var (
    barVal = map[string]int{"alpha": 34, "bravo": 56, "charlie": 23,
                            "delta": 87, "echo": 56, "foxtrot": 12,
                            "golf": 34, "hotel": 16, "indio": 87,
                            "juliet": 65, "kili": 43, "lima": 98}
)

func main() {
    fmt.Println("unsorted:")
    for k, v := range barVal {
        fmt.Printf("Key: %v, Value: %v / ", k, v)
    }
    keys := make([]string, len(barVal))
    i := 0
    for k, _ := range barVal {
        keys[i] = k
        i++
    }
    sort.Strings(keys)
    fmt.Println()
    fmt.Println("sorted:")
    for _, k := range keys {
        fmt.Printf("Key: %v, Value: %v / ", k, barVal[k])
    }
}
```

2. 但是若想要一个排序好的列表，还是使用结构体切片会比较有效：

```go
type name struct {
    key string
    value int
}
```
</details>

## 包(Package)

[包列表查询](https://gowalker.org/search?q=gorepos)

这一章主要讲各种库，以及自编库和编译安装到注意事项，故无太多记录。

安装外部库的命令为```go install xxx.com/xxx/yyy```(类似这样的，不一定是网址类型)

## 结构和方法

Go中没有类，所以struct的概念相比其他的语言来讲会更重要一些

<details>
    <summary>定义</summary>

```go
type identifier struct {
    field1 type1
    field2 type2
    ...
}

// type 1
var s T
s.a = 5
s.b = 8

// type 2
var t *T
t = new(T)
```

通过结构体的两种类型声明而出的一个是实例（指针变量）一个是对象；当给结构体别名的时候，两种类型可以互相直接转换
</details>

<details>
    <summary>通过工厂方法实现类似其他语言的构造函数</summary>

```go
type File struct {
    fd      int     // 文件描述符
    name    string  // 文件名
}
func NewFile(fd int, name string) *File {
    if fd < 0 {
        return nil
    }

    return &File{fd, name}
}
f := NewFile(10, "./test.txt")
```

强制使用工厂方法：只需要将包的结构体用小写开头，其他包则无法直接访问到该类型，只能通过可见的工厂方法来构造这个实例。
</details>

<details>
    <summary>带标签的结构体</summary>

```go
type TagType struct { // tags
    field1 bool   "An important answer"
    field2 string "The name of the thing"
    field3 int    "How much there are"
}
// 其中的field类型后的字符串就是tag，可以通过反射来获取类型，然后通过下标获取字段，通过字段的Tag属性来获取这个字符串。
```
</details>

<details>
    <summary>匿名字段和内嵌结构体</summary>

结构体中可以内嵌有类型的而无变量名的结构体变量，然后可以直接获取到相应变量中的字段等，内嵌变量（如int，float也是可以的）

```go
package main

import "fmt"

type innerS struct {
    in1 int
    in2 int
}

type outerS struct {
    b    int
    c    float32
    int  // anonymous field
    innerS //anonymous field
}

func main() {
    outer := new(outerS)
    outer.b = 6
    outer.c = 7.5
    outer.int = 60
    outer.in1 = 5
    outer.in2 = 10

    fmt.Printf("outer.b is: %d\n", outer.b)
    fmt.Printf("outer.c is: %f\n", outer.c)
    fmt.Printf("outer.int is: %d\n", outer.int)
    fmt.Printf("outer.in1 is: %d\n", outer.in1)
    fmt.Printf("outer.in2 is: %d\n", outer.in2)

    // 使用结构体字面量
    outer2 := outerS{6, 7.5, 60, innerS{5, 10}}
    fmt.Println("outer2 is:", outer2)
}

// 输出：
// outer.b is: 6
// outer.c is: 7.500000
// outer.int is: 60
// outer.in1 is: 5
// outer.in2 is: 10
// outer2 is:{6 7.5 60 {5 10}}
```

当命名冲突(内嵌不同结构体中的变量名重复)的时候，外部覆盖内部，如果处于同一层，需要程序员明确指定是哪个类型中的属性
</details>

<details>
    <summary>方法</summary>

结构体+方法近似于OO中的类。方法是有接收者的函数，声明方法如下：

```go
func (recv receiver_type) methodName(parameter_list) (return_value_list) { ... }
```

1. receiver_type可以为任意类型（在相同包中声明），但是不能为接口、指针类型（但是可以是允许的类型的指针）
2. 当接收者是指针的时候，可以在方法中修改接收者的值或者状态
3. 指针方法和值方法都可以在指针或非指针上被调用，如下面程序所示，类型 List 在值上有一个方法 Len()，在指针上有一个方法 Append()，但是可以看到两个方法都可以在两种类型的变量上被调用。

```go
package main

import (
    "fmt"
)

type List []int

func (l List) Len() int        { return len(l) }
func (l *List) Append(val int) { *l = append(*l, val) }

func main() {
    // 值
    var lst List
    lst.Append(1)
    fmt.Printf("%v (len: %d)", lst, lst.Len()) // [1] (len: 1)

    // 指针
    plst := new(List)
    plst.Append(2)
    fmt.Printf("%v (len: %d)", plst, plst.Len()) // &[2] (len: 1)
}
```
</details>


<details>
    <summary>String()方法</summary>

通过定义类型的String方法，当调用```fmt.Println(struct_obj)```的时候，会输出String中的方法，调试方便。
</details>

## 接口





