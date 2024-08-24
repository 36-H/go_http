package geecache

import (
	"container/list"
)

type Fifo struct {
	// 缓存最大的容量，单位字节；
	maxBytes int
	// 当一个 entry 从缓存中移除是调用该回调函数，默认为 nil
	onEvicted func(key string, value interface{})

	// 已使用的字节数，只包括值，key 不算
	usedBytes int

	ll    *list.List
	cache map[string]*list.Element
}

// New 创建一个新的 Cache，如果 maxBytes 是 0，表示没有容量限制
func NewFIFO(maxBytes int, onEvicted func(key string, value interface{})) cacheType {
	return &Fifo{
		maxBytes:  maxBytes,
		onEvicted: onEvicted,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
	}
}

// 增改
func (fifo *Fifo) Put(key string, value interface{}) {
	if e, ok := fifo.cache[key]; ok {
		fifo.ll.MoveToBack(e)
		en := e.Value.(*entry)
		fifo.usedBytes = fifo.usedBytes - Len(en.Value) + Len(value)
		en.Value = value
		return
	}

	en := &entry{Key: key, Value: value}
	e := fifo.ll.PushBack(en)
	fifo.cache[key] = e
	// fmt.Printf("当前使用容量:%d,",fifo.usedBytes)
	// fmt.Printf("即将加入缓存的字节数:%d,缓存类型:%T,",Len(en.Value),value)
	fifo.usedBytes += Len(en.Value)
	// fmt.Printf("新使用容量:%d\n",fifo.usedBytes)
	for fifo.maxBytes > 0 && fifo.usedBytes > fifo.maxBytes {
		fifo.RemoveOldest()
	}
}

// 查
func (fifo *Fifo) Get(key string) (interface{}, bool) {
	if e, ok := fifo.cache[key]; ok {
		return e.Value.(*entry).Value,true
	}

	return nil,false
}

// 删
func (fifo *Fifo) Remove(key string) {
	if e, ok := fifo.cache[key]; ok {
		fifo.removeElement(e)
	}
}

// 淘汰
func (fifo *Fifo) RemoveOldest() {
	fifo.removeElement(fifo.ll.Front())
}

// 长度
func (fifo *Fifo) Len() int {
	return fifo.ll.Len()
}

// 当前容量
func (fifo *Fifo) Size() int {
	return fifo.usedBytes
}

// 总容量
func (fifo *Fifo) Capacity() int {
	return fifo.maxBytes
}

func (fifo *Fifo) removeElement(e *list.Element) {
	if e == nil {
		return
	}

	fifo.ll.Remove(e)
	en := e.Value.(*entry)
	fifo.usedBytes -= Len(en.Value)
	delete(fifo.cache, en.Key)

	if fifo.onEvicted != nil {
		fifo.onEvicted(en.Key, en.Value)
	}
}
