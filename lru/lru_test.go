package lru

import (
	"reflect"
	"testing"
)

type String string

func (str String) Len() uint {
	return uint(len(str))
}

func TestCache_Get(t *testing.T) {
	cache := New(0, nil)
	cache.Add("name", String("1234"))
	if val, ok := cache.Get("name"); !ok && val != String("1234") {
		t.Fatalf("cache hit name=1234 fail")
	}
	if _, ok := cache.Get("key1"); ok {
		t.Fatalf("cache has key1, fail")
	}
}

func TestCache_Remote(t *testing.T) {
	k1, k2, k3 := "k1", "k2", "k3"
	v1, v2, v3 := "v1", "v2", "v3"
	caps := uint(len(k1 + k2 + v1 + v2))

	cache := New(caps, nil)
	cache.Add(k1, String(v1))
	cache.Add(k2, String(v2))
	cache.Add(k3, String(v3))

	if _, ok := cache.Get(k1); ok || cache.Len() != 2 {
		t.Fatalf("Removeoldest key1 failed")
	}
}

func TestCache_Evicted(t *testing.T) {
	keys := make([]string, 0)

	callback := func(key string, value Value) {
		keys = append(keys, key)
	}
	cache := New(10, callback)

	cache.Add("key1", String("12345"))
	cache.Add("key2", String("v2"))
	cache.Add("key3", String("v3"))
	cache.Add("key4", String("v4"))

	expect := []string{"key1", "key2", "key3"}

	if !reflect.DeepEqual(keys, expect) {
		t.Fatalf("value is not equal")
	}

}
