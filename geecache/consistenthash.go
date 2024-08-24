package geecache

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// Hash maps bytes to uint32
type Hash func(data []byte) uint32

// Map constains all hashed keys
type Map struct {
	//hash 函数
	hash     Hash
	//虚拟节点倍数
	replicas int
	//哈希环
	keys     []int // Sorted
	//虚拟节点与真实节点的映射表 hashMap，键是虚拟节点的哈希值，值是真实节点的名称。
	hashMap  map[int]string
}

// 创建一致性Hash
func NewConsistentHash(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// 增加节点
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
}

// Get 获取哈希环中与所提供的键最接近的项
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	//计算 key 的哈希值
	hash := int(m.hash([]byte(key)))
	//顺时针找到第一个匹配的虚拟节点的下标 idx，从 m.keys 中获取到对应的哈希值。
	//如果 idx == len(m.keys)，说明应选择 m.keys[0]，
	//因为 m.keys 是一个环状结构，所以用取余数的方式来处理这种情况。
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	//通过 hashMap 映射得到真实的节点。
	return m.hashMap[m.keys[idx%len(m.keys)]]
}