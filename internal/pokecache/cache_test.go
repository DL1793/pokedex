package pokecache

import (
	"fmt"
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("This is a test"),
		},
		{
			key: "https://example.com/path",
			val: []byte("More test data"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := NewCache(interval)
			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("Could not find key %v", c.key)
				return
			}
			if string(val) != string(c.val) {
				t.Errorf("Got %v, expected %v", string(val), string(c.val))
				return
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime + 5*time.Millisecond
	cache := NewCache(baseTime)
	cache.Add("https://example.com", []byte("This is a test"))
	_, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("Expected to find a key")
		return
	}
	time.Sleep(waitTime)
	_, ok = cache.Get("https://example.com")
	if ok {
		t.Errorf("Expected to not find a key")
		return
	}
}
