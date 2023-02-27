package cache

import (
	"fmt"
	"log"
	"sync"
)

type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (get GetterFunc) Get(key string) ([]byte, error) {
	return get(key)
}

type Group struct {
	name      string
	getter    Getter
	mainCache cache
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, maxBytes uint, getter Getter) (g *Group) {
	if getter == nil {
		panic("Getter is nil")
	}
	mu.Lock()
	defer mu.Unlock()

	g = &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: maxBytes},
	}

	groups[name] = g
	return
}

func GetGroup(name string) *Group {
	// 读写锁
	mu.RLock()
	defer mu.RUnlock()
	return groups[name]
}

func (g *Group) populateCache(key string, val ByteView) {
	g.mainCache.add(key, val)
}

func (g *Group) getLocally(key string) (val ByteView, err error) {
	bytes, err := g.getter.Get(key)

	if err != nil {
		val = ByteView{}
	} else {
		val = ByteView{b: cloneBytes(bytes)}
		g.populateCache(key, val)
	}
	return
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}
	if val, ok := g.mainCache.get(key); ok {
		log.Println("[GeeCache] hit")
		return val, nil
	}
	return g.getLocally(key)
}
