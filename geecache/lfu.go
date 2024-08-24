package geecache

import (
	"container/heap"
)

type Lfu struct {
	// 缓存最大的容量，单位字节；
	maxBytes int
	// 当一个 entry 从缓存中移除是调用该回调函数，默认为 nil
	onEvicted func(key string, value interface{})

	// 已使用的字节数，只包括值，key 不算
	usedBytes int

	queue *queue
	cache map[string]*lfu_entry
}

// New 创建一个新的 Cache，如果 maxBytes 是 0，表示没有容量限制
func NewLFU(maxBytes int, onEvicted func(key string, value interface{})) cacheType {
	q := make(queue, 0, 1024)
	return &Lfu{
		maxBytes:  maxBytes,
		onEvicted: onEvicted,
		queue:     &q,
		cache:     make(map[string]*lfu_entry),
	}
}

// Set 往 Cache 增加一个元素（如果已经存在，更新值，并增加权重，重新构建堆）
func (lfu *Lfu) Put(key string, value interface{}) {
	if e, ok := lfu.cache[key]; ok {
		lfu.usedBytes = lfu.usedBytes - Len(e.Value) + Len(value)
		lfu.queue.update(e, value, e.weight+1)
		return
	}

	en := &lfu_entry{
		entry: entry{
			Key:   key,
			Value: value,
		},
	}
	heap.Push(lfu.queue, en)
	lfu.cache[key] = en
	lfu.usedBytes += Len(en.Value)
	if lfu.maxBytes > 0 && lfu.usedBytes > lfu.maxBytes {
		lfu.removeElement(heap.Pop(lfu.queue))
	}
}

// Get 从 cache 中获取 key 对应的值，nil 表示 key 不存在
func (lfu *Lfu) Get(key string) (interface{}, bool) {
	if e, ok := lfu.cache[key]; ok {
		lfu.queue.update(e, e.Value, e.weight+1)
		return e.Value, true
	}

	return nil, false
}

// Remove 从 cache 中删除 key 对应的元素
func (lfu *Lfu) Remove(key string) {
	if e, ok := lfu.cache[key]; ok {
		heap.Remove(lfu.queue, e.index)
		lfu.removeElement(e)
	}
}

// RemoveOldest 从 cache 中删除最旧的记录
func (lfu *Lfu) RemoveOldest() {
	if lfu.queue.Len() == 0 {
		return
	}
	lfu.removeElement(heap.Pop(lfu.queue))
}

func (lfu *Lfu) removeElement(e interface{}) {
	if e == nil {
		return
	}

	en := e.(*lfu_entry)

	delete(lfu.cache, en.Key)

	lfu.usedBytes -= Len(en.Value)

	if lfu.onEvicted != nil {
		lfu.onEvicted(en.Key, en.Value)
	}
}

// Len 返回当前 cache 中的记录数
func (lfu *Lfu) Len() int {
	return lfu.queue.Len()
}

// 当前容量
func (lfu *Lfu) Size() int {
	return lfu.usedBytes
}

// 总容量
func (lfu *Lfu) Capacity() int {
	return lfu.maxBytes
}
