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
- [包(Package)](#包)
- [结构(struct)和方法(method)](#结构和方法)
- [接口](#接口)
- [读写](#读写)
- [错误处理和测试](#错误处理和测试)
- [协程与通道](#协程与通道)
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

> 就是自己的编译器可以自行编译自己的编译器。  
> 实现方法就是这个编译器的作者用这个语言的一些特性来编写编译器并在该编译器中支持这些自己使用到的特性。  
> 首先，第一个编译器肯定是用别的语言写的（不论是C还是Go还是Lisp还是Python），后面的版本才能谈及自举。  
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

[Interface相关](#接口)

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

## 包

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

<details>
    <summary>接口定义</summary>

一组方法的集合，不包含实现代码。

1. 类型不需要显式声明它实现了某个接口：接口被隐式地实现。多个类型可以实现同一个接口。
2. 实现某个接口的类型（除了实现接口方法外）可以有其他的方法。
3. 一个类型可以实现多个接口。
4. 接口类型可以包含一个实例的引用，该实例的类型实现了此接口（接口是动态类型）。
5. 接口可以嵌套接口。

Example：

```go
package main

import "fmt"

type Shaper interface {
    Area() float32
}

type Square struct {
    side float32
}

func (sq *Square) Area() float32 {
    return sq.side * sq.side
}

type Rectangle struct {
    length, width float32
}

func (r Rectangle) Area() float32 {
    return r.length * r.width
}

func main() {

    r := Rectangle{5, 3} // Area() of Rectangle needs a value
    q := &Square{5}      // Area() of Square needs a pointer
    // shapes := []Shaper{Shaper(r), Shaper(q)}
    // or shorter
    shapes := []Shaper{r, q}
    fmt.Println("Looping through shapes for area ...")
    for n, _ := range shapes {
        fmt.Println("Shape details: ", shapes[n])
        fmt.Println("Area of this shape is: ", shapes[n].Area())
    }
}

// Looping through shapes for area ...
// Shape details:  {5 3}
// Area of this shape is:  15
// Shape details:  &{5}
// Area of this shape is:  25
```
</details>

<details>
    <summary>检测接口变量的类型</summary>

接口变量可以包含实例的引用，而很多时候我们需要确定该引用类型。我们假定接口变量名为var_inter，类型名为struct_type，那么我们可以通过：

```go
val,ok := var_inter.(*struct_type)
// ok为true时，val是转换后的值；否则为该类型空值

// 另一种方式判断类型
switch t := var_inter.(type) {
case *Square:
    fmt.Printf("Type Square %T with value %v\n", t, t)
case *Circle:
    fmt.Printf("Type Circle %T with value %v\n", t, t)
case nil:
    fmt.Printf("nil value: nothing to check?\n")
default:
    fmt.Printf("Unexpected type %T\n", t)
}
// Type Square *main.Square with value &{5}
```
</details>

<details>
    <summary>确定某类型是否实现了某接口</summary>

```go
// 假定v是一个值
type Stringer interface {
    String() string
}

if sv, ok := v.(Stringer); ok {
    fmt.Printf("v implements String(): %s\n", sv.String()) // note: sv, not v
}
```
</details>

<details>
    <summary>接口接收者可调用类型</summary>

1. 指针方法可以通过指针调用
2. 值方法可以通过值调用
3. 接收者是值的方法可以通过指针调用，因为指针会首先被解引用
4. 接收者是指针的方法不可以通过值调用，因为存储在接口中的值没有地址

</details>

<details>
    <summary>空接口的作用</summary>

1. 由于空接口能承接任意类型的变量，所以可以实现承接任意类型的切片
2. 实现数据结构（如树）的时候，data字段可以用空接口，这样就能存储任意的值，使得代码具有足够的通用性
</details>

<details>
    <summary>通过反射修改值</summary>

[源码](./src/reflect.go)
</details>

<details>
    <summary>类型内嵌接口</summary>

> PS：这里原书那里的一些说法与我测试结果不太相同

当一个类型包含（内嵌）另一个类型（实现了一个或多个接口）的**指针**时，这个类型就可以使用（另一个类型）所有的接口方法。更无歧义的表达是：

1. 接口可以内嵌接口
2. 结构体可以内嵌结构体或结构体指针
3. 结构体可以内嵌接口，此时初始化时要用实现了该接口的类型来初始化

[代码](./src/embeded_Interface.go)

[有关内嵌类型的阅读](https://travix.io/type-embedding-in-go-ba40dd4264df)

接口可以通过继承多个接口来提供像**多重继承**一样的特性
</details>

## 读写

<details>
    <summary>读取键盘输入</summary>

fmt包提供了Scan或Sscan开头的函数（Scanln和Sscanf），其中Scanln以空格分隔符，直到遇到换行；Sscanf则类似c中的scanf，按照第一个参数规定的顺序来获取输入。

<details>
        <summary>Example:</summary>

```go
// 从控制台读取输入:
package main
import "fmt"

var (
   firstName, lastName, s string
   i int
   f float32
   input = "56.12 / 5212 / Go"
   format = "%f / %d / %s"
)

func main() {
   fmt.Println("Please enter your full name: ")
   fmt.Scanln(&firstName, &lastName)
   // fmt.Scanf("%s %s", &firstName, &lastName)
   fmt.Printf("Hi %s %s!\n", firstName, lastName) // Hi Chris Naegels
   fmt.Sscanf(input, format, &f, &i, &s)
   fmt.Println("From the string we read: ", f, i, s)
    // 输出结果: From the string we read: 56.12 5212 Go
}
```
</details>

也可以使用 bufio 包提供的缓冲读取（buffered reader）来读取数据：

<details>
    <summary>Example:</summary>

```go
package main
import (
    "fmt"
    "bufio"
    "os"
)

var inputReader *bufio.Reader
var input string
var err error

func main() {
    inputReader = bufio.NewReader(os.Stdin)
    fmt.Println("Please enter some input: ")
    input, err = inputReader.ReadString('\n')
    if err == nil {
        fmt.Printf("The input was: %s\n", input)
    }
}
```
</details>

</details>

<details>
    <summary>文件读</summary>

os.File类型的指针表示文件句柄，os.Stdin和os.Stdout的类型都是\*os.File

<details>
    <summary>文件读示例</summary>

```go
package main

import (
    "bufio"
    "fmt"
    "io"
    "os"
)

func main() {
    inputFile, inputError := os.Open("input.dat")
    if inputError != nil {
        fmt.Printf("An error occurred on opening the inputfile\n" +
            "Does the file exist?\n" +
            "Have you got acces to it?\n")
        return // exit the function on error
    }
    defer inputFile.Close()

    inputReader := bufio.NewReader(inputFile)
    for {
        inputString, readerError := inputReader.ReadString('\n')
        fmt.Printf("The input was: %s", inputString)
        if readerError == io.EOF {
            return
        }      
    }
}
```
</details>

而带缓冲的读取二进制文件的方法，可以用Read()函数来处理

```go
buf := make([]byte, 1024)
for {
    n, err := inputReader.Read(buf)
    if (n == 0) {
        break
    }
}
```

[完整示例](./src/file_read.go)


压缩文件的读取，利用compress包读取

</details>


<details>
    <summary>文件写</summary>

```go
package main

import "os"

func main() {
    os.Stdout.WriteString("hello, world\n")
    f, _ := os.OpenFile("test", os.O_CREATE|os.O_WRONLY, 0666)
    defer f.Close()
    f.WriteString("hello, world in a file\n")
}
```

<details>
    <summary>使用bufio的方式</summary>

```go
package main

import (
    "os"
    "bufio"
    "fmt"
)

func main () {
    // var outputWriter *bufio.Writer
    // var outputFile *os.File
    // var outputError os.Error
    // var outputString string
    outputFile, outputError := os.OpenFile("output.dat", os.O_WRONLY|os.O_CREATE, 0666)
    if outputError != nil {
        fmt.Printf("An error occurred with file opening or creation\n")
        return  
    }
    defer outputFile.Close()

    outputWriter := bufio.NewWriter(outputFile)
    outputString := "hello world!\n"

    for i:=0; i<10; i++ {
        outputWriter.WriteString(outputString)
    }
    outputWriter.Flush()
}
```
</details>

</details>

<details>
    <summary>JSON的序列化和反序列化</summary>

json的库在```encoding/json```，其中序列化函数```json.Marshal()```的函数签名是```func Marshal(v interface{}) ([]byte, error)```，反序列化的函数```UnMarshal()```的函数签名是```func Unmarshal(data []byte, v interface{}) error```

解码的时候要注意解码后格式的转换。

<details>
    <summary>例子</summary>

```go
b := []byte(`{"Name": "Wednesday", "Age": 6, "Parents": ["Gomez", "Morticia"]}`)
var f interface{}
err := json.Unmarshal(b, &f)

//  f指向的值是一个 map，key 是一个字符串，value 是自身存储作为空接口类型的值：
map[string]interface{} {
    "Name": "Wednesday",
    "Age":  6,
    "Parents": []interface{} {
        "Gomez",
        "Morticia",
    },
}

// 要访问这个数据，我们可以使用类型断言
m := f.(map[string]interface{})

// 我们可以通过 for range 语法和 type switch 来访问其实际类型：
for k, v := range m {
    switch vv := v.(type) {
    case string:
        fmt.Println(k, "is string", vv)
    case int:
        fmt.Println(k, "is int", vv)

    case []interface{}:
        fmt.Println(k, "is an array:")
        for i, u := range vv {
            fmt.Println(i, u)
        }
    default:
        fmt.Println(k, "is of a type I don’t know how to handle")
    }
}
// 通过这种方式，你可以处理未知的 JSON 数据，同时可以确保类型安全。


```
</details>
</details>


## 错误处理和测试

Go中处理普通错误，应该通过函数最后一个返回值返回个主调方，如果返回nil表示正常。至于panic and recover是用在真正的异常上的（无法预测的错误上的）。fmt中也有```fmt.Errorf()```来打印错误信息，使用方法与```fmt.Printf()```一模一样。


<details>
    <summary>自定义错误</summary>

```go
// PathError records an error and the operation and file path that caused it.
type PathError struct {
    Op string    // "open", "unlink", etc.
    Path string  // The associated file.
    Err error  // Returned by the system call.
}

func (e *PathError) String() string {
    return e.Op + " " + e.Path + ": "+ e.Err.Error()
}
```
</details>

<details>
    <summary>panic和recover</summary>

recover只能在defer修饰的函数中使用：用于取得panic调用中传递过来的错误值，如果是正常执行，调用recover会返回nil，且没有其它效果。

一个简单的例子：

```go
// panic_recover.go
package main

import (
    "fmt"
)

func badCall() {
    panic("bad end")
}

func test() {
    defer func() {
        if e := recover(); e != nil {
            fmt.Printf("Panicing %s\r\n", e)
        }
    }()
    badCall()
    fmt.Printf("After bad call\r\n") // <-- wordt niet bereikt
}

func main() {
    fmt.Printf("Calling test\r\n")
    test()
    fmt.Printf("Test completed\r\n")
}

// Calling test
// Panicing bad end
// Test completed
``` 

> 计算机科学领域的任何问题都可以通过增加一个简介的中间层来解决。  
> Any problem in computer science can be solved by another layer of indirection.

所以用以下闭包的方式（外层包装一个error_handler，并于其中的defer进行recover）来解决多次判断错误的不优雅代码：[传送门](https://go.fdos.me/13.5.html)

</details>

<details>
    <summary>测试</summary>

测试代码写于xx_test.go中，即当源码文件为add.go的时候，测试代码为add_test.go。且测试数据通常通过表驱动的方式，在函数中for循环对比输入输出是否正确。
    
</details>

## 协程与通道

通过go func()关键字来调用goroutine

<details>
    <summary>goroutine1.go</summary>

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    fmt.Println("In main()")
    go longWait()
    go shortWait()
    fmt.Println("About to sleep in main()")
    // sleep works with a Duration in nanoseconds (ns) !
    time.Sleep(10 * 1e9)
    fmt.Println("At the end of main()")
}

func longWait() {
    fmt.Println("Beginning longWait()")
    time.Sleep(5 * 1e9) // sleep for 5 seconds
    fmt.Println("End of longWait()")
}

func shortWait() {
    fmt.Println("Beginning shortWait()")
    time.Sleep(2 * 1e9) // sleep for 2 seconds
    fmt.Println("End of shortWait()")
}

// In main()
// About to sleep in main()
// Beginning longWait()
// Beginning shortWait()
// End of shortWait()
// End of longWait()
// At the end of main()
```
</details>

<details>
    <summary>通道</summary>

声明方式：
> var identifier chan datatype

```go
var ch1 chan string
ch1 = make(chan string)

// or

ch1 := make(chan string)

// 发送数据
var text string
text = "Hello world" 
ch <- text

// 接收数据
output := <- ch
```

<details>
    <summary>goroutine2.go</summary>

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    ch := make(chan string)

    go sendData(ch)
    go getData(ch)

    time.Sleep(1e9)
}

func sendData(ch chan string) {
    ch <- "Washington"
    ch <- "Tripoli"
    ch <- "London"
    ch <- "Beijing"
    ch <- "Tokio"
}

func getData(ch chan string) {
    var input string
    // time.Sleep(2e9)
    for {
        input = <-ch
        fmt.Printf("%s ", input)
    }
}

// 输出如下：
// Washington Tripoli London Beijing Tokio
```

</details>

容量为0的通道是阻塞的，即发送和接受操作都是阻塞的，发送者或接收者未就绪的时候，通道都是阻塞的，通道使用中，对于新的输入也是阻塞的。

声明带缓冲的通道（异步的非阻塞，满或空的时候还是阻塞的）方法：```ch := make(chan type, value)```

<details>
    <summary>排序中使用通道</summary>

```go
done := make(chan bool)
// doSort is a lambda function, so a closure which knows the channel done:
doSort := func(s []int){
    sort(s)
    done <- true
}
i := pivot(s)
go doSort(s[:i])
go doSort(s[i:])
<-done
<-done
```
</details>


</details>




