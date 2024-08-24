package geecache

import (
	"fmt"
	"log"
	"sync"
)

// 加载键的数据
type Getter interface {
	Get(key string) ([]byte, error)
}

// A GetterFunc implements Getter with a function.
type GetterFunc func(key string) ([]byte, error)

// Get implements Getter interface function
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

// Group 是一个缓存命名空间，关联的数据加载分布
type Group struct {
	name      string
	getter    Getter
	mainCache Cache
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

// NewGroup 创建一个新的组实例
func NewGroup(name string, cacheBytes int, getter Getter,onEvicted func(key string, value interface{}),t T) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: *NewCache(cacheBytes,onEvicted,t),
	}
	groups[name] = g
	return g
}

// GetGroup 返回先前使用 NewGroup 创建的命名组，
// 或者如果没有这样的组则为零。
func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}


// Get 从缓存中获取键的值
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	if v, ok := g.mainCache.get(key); ok {
		log.Println("[GeeCache] hit")
		return v, nil
	}

	return g.load(key)
}

func (g *Group) load(key string) (value ByteView, err error) {
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err

	}
	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
