# Channel源码学习

goversion：1.11.2

<details>
<summary>点击展开目录菜单</summary>
<!-- TOC -->

- [数据结构](#数据结构)
- [channel的重要函数](#channel的重要函数)
  - [makechan](#makechan)

<!-- /TOC -->

</details>

## 数据结构

```go
type hchan struct {
	qcount   uint           // total data in the queue
	dataqsiz uint           // size of the circular queue
	buf      unsafe.Pointer // points to an array of dataqsiz elements
	elemsize uint16
	closed   uint32
	elemtype *_type // element type
	sendx    uint   // send index
	recvx    uint   // receive index
	recvq    waitq  // list of recv waiters
	sendq    waitq  // list of send waiters

	// lock protects all fields in hchan, as well as several
	// fields in sudogs blocked on this channel.
	//
	// Do not change another G's status while holding this lock
	// (in particular, do not ready a G), as this can deadlock
	// with stack shrinking.
	lock mutex
}
type waitq struct {
	first *sudog
	last  *sudog
}
type sudog struct {
	// The following fields are protected by the hchan.lock of the
	// channel this sudog is blocking on. shrinkstack depends on
	// this for sudogs involved in channel ops.

	g *g

	// isSelect indicates g is participating in a select, so
	// g.selectDone must be CAS'd to win the wake-up race.
	isSelect bool
	next     *sudog
	prev     *sudog
	elem     unsafe.Pointer // data element (may point to stack)

	// The following fields are never accessed concurrently.
	// For channels, waitlink is only accessed by g.
	// For semaphores, all fields (including the ones above)
	// are only accessed when holding a semaRoot lock.

	acquiretime int64
	releasetime int64
	ticket      uint32
	parent      *sudog // semaRoot binary tree
	waitlink    *sudog // g.waiting list or semaRoot
	waittail    *sudog // semaRoot
	c           *hchan // channel
}
```

从这些数据结构可以看出sudog在g(goroutine)的基础上做了一定的封装，用于记录一个在channel上等待的协程g和等待的元素elem。

## channel的重要函数

### makechan

<details>
	<summary>展开</summary>

```go
func makechan(t *chantype, size int) *hchan {
	elem := t.elem

	// compiler checks this but be safe.
	if elem.size >= 1<<16 {
		throw("makechan: invalid channel element type")
	}
	if hchanSize%maxAlign != 0 || elem.align > maxAlign {
		throw("makechan: bad alignment")
	}

	if size < 0 || uintptr(size) > maxSliceCap(elem.size) || uintptr(size)*elem.size > maxAlloc-hchanSize {
		panic(plainError("makechan: size out of range"))
	}

	// Hchan does not contain pointers interesting for GC when elements stored in buf do not contain pointers.
	// buf points into the same allocation, elemtype is persistent.
	// SudoG's are referenced from their owning thread so they can't be collected.
	// TODO(dvyukov,rlh): Rethink when collector can move allocated objects.
	var c *hchan
	switch {
	case size == 0 || elem.size == 0:
		// Queue or element size is zero.
		c = (*hchan)(mallocgc(hchanSize, nil, true))
		// Race detector uses this location for synchronization.
		c.buf = c.raceaddr()
	case elem.kind&kindNoPointers != 0:
		// Elements do not contain pointers.
		// Allocate hchan and buf in one call.
		c = (*hchan)(mallocgc(hchanSize+uintptr(size)*elem.size, nil, true))
		c.buf = add(unsafe.Pointer(c), hchanSize)
	default:
		// Elements contain pointers.
		c = new(hchan)
		c.buf = mallocgc(uintptr(size)*elem.size, elem, true)
	}

	c.elemsize = uint16(elem.size)
	c.elemtype = elem
	c.dataqsiz = uint(size)

	if debugChan {
		print("makechan: chan=", c, "; elemsize=", elem.size, "; elemalg=", elem.alg, "; dataqsiz=", size, "\n")
	}
	return c
}
```

#### hchanSize简易证明

```makechan()```中的hchanSize的实现是一个很有意思的二进制trick，实现如下：

```go
const (
	maxAlign  = 8
	hchanSize = unsafe.Sizeof(hchan{}) + uintptr(-int(unsafe.Sizeof(hchan{}))&(maxAlign-1))
	debugChan = false
)
```

这段代码的意思是，将hchanSize设置为比hchan大的最小的maxAlign的倍数，至于为什么可以达到这个效果，这里有个简单的不算特别严谨的证明：

<details>
	<summary>证明展开</summary>

```a + ( ( -a ) & (alignSize - 1) )```能计算出大于等于a的最小的 ( alignSize的倍数 )，前提条件是alignSize是2的次幂

也就是```a + ( ( -a ) & (alignSize - 1) ) == a + alignSize - a % alignSize```

举个例子：a = 50, alignSize = 8 上式结果就是56（alignSize是2的n次幂）

具体证明：

已知 ```a % alignSize == a & ( alignSize - 1) ```，其中```alignSize = 2^n```

要证```a + ( ( -a ) & (alignSize - 1) ) == a + alignSize - a % alignSize```

即要证：``` ( -a ) & (alignSize - 1) == alignSize - ( a & (alignSize-1) ) ```

即要证：``` ( -a ) & (alignSize - 1) + ( a & (alignSize - 1) ) == alignSize ```

即要证：``` ( -a ) & (2^n - 1) + ( a & (2^n - 1) ) == 2^n ```

即要证：``` ( ～a + 1 ) & (2^n - 1) + ( a & (2^n - 1) ) == 2^n ```

由于alignSize为2^n，那么```x & (2^n-1)```其实就是将高于2^n的高位截断

而-a的计算机实现是原数取反并二进制+1，那么在可表示范围内截断，相加其实就是等于这个二的次幂的。

</details>

#### makechan正文

前三个if是用于检查的，第一个检查channel元素大小限制（之所以左移16位估计是因为后续hchan中的elemsize类型是uint16）；第二个检查是检查内存对齐的，这个跟上面的计算hchanSize有关；第三个确认size的值没有溢出，且内存足够分配。

接下来是内存分配的代码

1. 如果size为0或者elem的size为0，那么只分配hchanSize的内存，c.buf实际上不会分配对应的内存

2. 如果类型不包含指针，则将channel和对应的elem所需的内存一次性连续分配

3. 如果类型内包含指针，则将两个内存的申请分开。下面说下这里为什么内存申请分开：

   之所以这么设计，是因为golang采用的三色标记法，golang中的内存分配时采用的如下结构管理span（用于内存管理的基本单位）

   ```go
   type mcache struct {
   	alloc [67*2]*mspan // 按class分组的mspan列表
   }
   
   ```

   这里面有个乘2，原因是分了两类span，一类是需要gc扫描的scan和不用gc扫描的noscan，noscan类型分配的内存中的对象是不包含指针类型的，所以三色标记法就不需要再去扫描他引用的对象，可以加快gc速度。再回到这个问题，看到makechan中有一行注释是这样的：```// Hchan does not contain pointers interesting for GC when elements stored in buf do not contain pointers.```

   可以看到其实hchan中并不会有三色标记法中的引用问题（hchan中的type类型仍然是指针，查阅网上资料貌似_type可以认为是Go中所有类型的公共描述，所以个人认为这个是程序运行过程中会持久引用的内容，不需要gc，所以hchan也可以视作没有指针的类型），那么只要elem里没有指针类型的时候，就可以放在一起申请内存，这样直接申请一块noscan类型的即可，在gc过程中就不需要再去扫描这一块的内存了。当elem中有指针的时候，hchan还是相当于没有指针的noscan类型，而elem此时是scan的类型，所以将这两个分开进行初始化。这就是这里分开初始化的原因。

然后后面就是继续将size等其他值赋值。

</details>

### chansend

<details> 
	<summary>展开</summary>

```go
/*
 * generic single channel send/recv
 * If block is not nil,
 * then the protocol will not
 * sleep but return if it could
 * not complete.
 *
 * sleep can wake up with g.param == nil
 * when a channel involved in the sleep has
 * been closed.  it is easiest to loop and re-run
 * the operation; we'll see that it's now closed.
 */
func chansend(c *hchan, ep unsafe.Pointer, block bool, callerpc uintptr) bool {...}
```

发送函数中的block主要是来自select的，如果一个channel来自于select，那么这个block参数会为false，这样就不会在select的时候被阻塞住。

代码中频繁用到的raceenabled常量，是用在go build -race时才会启用的，用于检测是否有race态，一般为false值，所以可以直接忽略这部分代码。

跟着继续往下走，如果channel已经关闭，那么直接Panic掉这个goroutine。

正常发送会有三种情况：

1. recvq队列中有空闲的接受者，那么直接发送给他
2. 如果队列未满，直接塞进队列中
3. 如果队列满了，且在select中，直接return false；否则阻塞队列

代码中也主要是这三段的不同处理，主要是代码中的一些细节的处理，第一种情况的话调用链是```send()```->```sendDirect()```代码如下：

```go
// runtime/chan.go#305
func sendDirect(t *_type, sg *sudog, src unsafe.Pointer) {
	// src is on our stack, dst is a slot on another stack.

	// Once we read sg.elem out of sg, it will no longer
	// be updated if the destination's stack gets copied (shrunk).
	// So make sure that no preemption points can happen between read & use.
	dst := sg.elem
    // 这里启用了内存屏障防止乱序执行后，导致先memmove，再执行dst := sg.elem
    // 使得一些操作修改了原来的elem而没有生效到dst上
	typeBitsBulkBarrier(t, uintptr(dst), uintptr(src), t.size)
	// No need for cgo write barrier checks because dst is always
	// Go memory.
	memmove(dst, src, t.size)
}
```

第二种情况中代码如下

```go
if c.qcount < c.dataqsiz {
    // Space is available in the channel buffer. Enqueue the element to send.
    qp := chanbuf(c, c.sendx)
    if raceenabled {
        // ...由于raceenabled在没指定情况下为false，忽略
    }
    typedmemmove(c.elemtype, qp, ep)
    c.sendx++
    if c.sendx == c.dataqsiz {
        c.sendx = 0
    }
    c.qcount++
    unlock(&c.lock)
    return true
}
```



</details>

