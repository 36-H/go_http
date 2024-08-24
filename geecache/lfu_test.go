package geecache

import (
	"testing"

	"github.com/matryer/is"
)

func TestSet(t *testing.T) {
	is := is.New(t)

	cache := NewLFU(24, nil)
	cache.RemoveOldest()
	cache.Put("k1", 1)
	v,_:= cache.Get("k1")
	is.Equal(v, 1)

	cache.Remove("k1")
	is.Equal(0, cache.Len())

	// cache.Set("k2", time.Now())
}

func TestLFUOnEvicted(t *testing.T) {
	is := is.New(t)

	keys := make([]string, 0, 8)
	onEvicted := func(key string, value interface{}) {
		keys = append(keys, key)
	}
	cache := NewLFU(16, onEvicted)

	cache.Put("k1", 1)
	cache.Put("k2", 2)
	// cache.Get("k1")
	// cache.Get("k1")
	// cache.Get("k2")
	cache.Put("k3", 3)
	cache.Put("k4", 4)

	expected := []string{"k1", "k3"}

	is.Equal(expected, keys)
	is.Equal(2, cache.Len())
}