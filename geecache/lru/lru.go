package lru

import (
	"container/list"
	"core"
)

type Lru struct {
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
func New(maxBytes int, onEvicted func(key string, value interface{})) core.CacheType {
	return &Lru{
		maxBytes:  maxBytes,
		onEvicted: onEvicted,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
	}
}

// 增改
func (lru *Lru) Put(key string, value interface{}) {
	if e, ok := lru.cache[key]; ok {
		lru.ll.MoveToBack(e)
		en := e.Value.(*core.Entry)
		lru.usedBytes = lru.usedBytes - core.Len(en.Value) + core.Len(value)
		en.Value = value
		return
	}

	en := &core.Entry{Key: key, Value: value}
	e := lru.ll.PushBack(en)
	lru.cache[key] = e
	// fmt.Printf("当前使用容量:%d,",Lru.usedBytes)
	// fmt.Printf("即将加入缓存的字节数:%d,缓存类型:%T,",core.Len(en.Value),value)
	lru.usedBytes += core.Len(en.Value)
	// fmt.Printf("新使用容量:%d\n",Lru.usedBytes)
	for lru.maxBytes > 0 && lru.usedBytes > lru.maxBytes {
		lru.RemoveOldest()
	}
}

// 查
func (lru *Lru) Get(key string) interface{} {
	if e, ok := lru.cache[key]; ok {
		lru.ll.MoveToBack(e)
		return e.Value.(*core.Entry).Value
	}

	return nil
}

// 删
func (lru *Lru) Remove(key string) {
	if e, ok := lru.cache[key]; ok {
		lru.removeElement(e)
	}
}

// 淘汰
func (lru *Lru) RemoveOldest() {
	lru.removeElement(lru.ll.Front())
}

// 长度
func (lru *Lru) Len() int {
	return lru.ll.Len()
}

// 当前容量
func (lru *Lru) Size() int {
	return lru.usedBytes
}

// 总容量
func (lru *Lru) Capacity() int {
	return lru.maxBytes
}

func (lru *Lru) removeElement(e *list.Element) {
	if e == nil {
		return
	}

	lru.ll.Remove(e)
	en := e.Value.(*core.Entry)
	lru.usedBytes -= core.Len(en.Value)
	delete(lru.cache, en.Key)

	if lru.onEvicted != nil {
		lru.onEvicted(en.Key, en.Value)
	}
}
