package geecache

import (
	"fmt"
	"runtime"
)

type T string

const (
	FIFO T = "fifo"
	LFU  T = "lfu"
	LRU  T = "lru"
)

func (t T) String() string {
	switch t {
	case FIFO:
		return "fifo"
	case LFU:
		return "lfu"
	case LRU:
		return "lru"
	default:
		return "lru"
	}

}

type cacheType interface {
	//增改
	Put(key string, value interface{})
	//查
	Get(key string) (interface{},bool)
	//删
	Remove(key string)
	//淘汰
	RemoveOldest()
	//当前容量
	Size() int
	//总容量
	Capacity() int
	//缓存个数
	Len() int
}

type entry struct {
	Key   string
	Value interface{}
}

type value interface {
	Len() int
}

func Len(val interface{}) int {
	// 不能使用如下方法统计 他会统计成interface{} 16字节
	// return int64(unsafe.Sizeof(value))
	var n int
	switch v := val.(type) {
	case value:
		n = v.Len()
	case string:
		if runtime.GOARCH == "amd64" {
			n = 16 + len(v)
		} else {
			n = 8 + len(v)
		}
	case bool, uint8, int8:
		n = 1
	case int16, uint16:
		n = 2
	case int32, uint32, float32:
		n = 4
	case int64, uint64, float64:
		n = 8
	case int, uint:
		if runtime.GOARCH == "amd64" {
			n = 8
		} else {
			n = 4
		}
	case complex64:
		n = 8
	case complex128:
		n = 16
	default:
		panic(fmt.Sprintf("%T is not implement value", val))
	}

	return n
}
