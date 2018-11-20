Go入门学习
===

目录
---
<details>
<summary>点击展开目录菜单</summary>

<!-- TOC -->

- [Go及环境相关](#Go及环境相关)
- [与其他语言不同之处](#与其他语言不同之处)

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