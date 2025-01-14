package geecache

import (
	"testing"

	"github.com/matryer/is"
)

func TestLRUOnEvicted(t *testing.T) {
	is := is.New(t)

	keys := make([]string, 0, 8)
	onEvicted := func(key string, value interface{}) {
		keys = append(keys, key)
	}
	cache := NewLRU(16, onEvicted)

	cache.Put("k1", 1)
	cache.Put("k2", 2)
	cache.Get("k1")
	cache.Put("k3", 3)
	cache.Get("k1")
	cache.Put("k4", 4)

	expected := []string{"k2", "k3"}

	is.Equal(expected, keys)
	is.Equal(2, cache.Len())
}