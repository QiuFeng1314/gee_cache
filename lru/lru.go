package lru

import (
	"container/list"
)

type Value interface {
	Len() uint
}

func (c *Cache) Len() uint {
	return uint(c.ll.Len())
}

type entry struct {
	key   string
	value Value
}

type Cache struct {
	maxBytes  uint       // 允许使用的最大内存
	nowBytes  uint       // 当前使用的内存
	ll        *list.List // 双向链表指针
	cache     map[string]*list.Element
	OnEvicted func(key string, value Value)
}

func New(maxBytes uint, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Get 查找
func (c *Cache) Get(key string) (Value, bool) {
	if element, ok := c.cache[key]; ok {
		c.ll.MoveToFront(element)
		kv := element.Value.(*entry)
		return kv.value, true
	}
	return nil, false
}

// Remote 缓存淘汰。移除最近最少访问的节点
func (c *Cache) Remote() {
	// 删除队首节点
	if back := c.ll.Back(); back != nil {
		c.ll.Remove(back)
		kv := back.Value.(*entry)
		// 删除Cache中的数据
		delete(c.cache, kv.key)
		// 更新当前使用的内存
		c.nowBytes -= uint(len(kv.key)) + kv.value.Len()
		// 执行回调
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Add 新增/修改
func (c *Cache) Add(key string, value Value) {
	if element, ok := c.cache[key]; ok {
		c.ll.MoveToFront(element)
		kv := element.Value.(*entry)
		kv.value = value
		c.nowBytes += uint(len(kv.key)) - kv.value.Len()
	} else {
		element := c.ll.PushFront(&entry{key: key, value: value})
		c.cache[key] = element
		c.nowBytes += uint(len(key)) + value.Len()
	}

	for c.maxBytes != 0 && c.maxBytes < c.nowBytes {
		c.Remote()
	}
}
