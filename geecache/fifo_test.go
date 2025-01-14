package geecache

import (
	"testing"

	"github.com/matryer/is"
)

func TestFIFOSetGet(t *testing.T) {
	is := is.New(t)

	cache := NewFIFO(24, nil)
	cache.RemoveOldest()
	cache.Put("k1", 1)
	v,_:= cache.Get("k1")
	is.Equal(v, 1)

	cache.Remove("k1")
	is.Equal(0, cache.Len()) // expect to be the same

	// cache.Set("k2", time.Now())
}

func TestOnEvicted(t *testing.T) {
	is := is.New(t)

	keys := make([]string, 0, 8)
	onEvicted := func(key string, value interface{}) {
		keys = append(keys, key)
	}
	cache := NewFIFO(16, onEvicted)

	cache.Put("k1", 1)
	cache.Put("k2", 2)
	cache.Get("k1")
	cache.Put("k3", 3)
	cache.Get("k1")
	cache.Put("k4", 4)

	expected := []string{"k1", "k2"}

	is.Equal(expected, keys)
	is.Equal(2, cache.Len())
}
