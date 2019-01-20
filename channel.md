# Channel源码学习

goversion：1.11.2

<details>
<summary>点击展开目录菜单</summary>
<!-- TOC -->

- [数据结构](#数据结构)
- [channel的重要函数](#channel的重要函数)

  - [makechan](#makechan)
  - 

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