package pokecache

import (
	"fmt"
	"sync"
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
			val: []byte("testdata"),
		},
		{
			key: "https://example.com/path",
			val: []byte("moretestdata"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := NewCache(interval)
			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("expected to find key")
				return
			}
			if string(val) != string(c.val) {
				t.Errorf("expected to find value")
				return
			}
		})
	}
}
func TestReapLoop(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime + 5*time.Millisecond
	cache := NewCache(baseTime)
	cache.Add("https://example.com", []byte("testdata"))

	_, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to find key")
		return
	}

	time.Sleep(waitTime)

	_, ok = cache.Get("https://example.com")
	if ok {
		t.Errorf("expected to not find key")
		return
	}
}

func TestMuLocking(t *testing.T) {
	initialTime := time.Now()
	const waitTime = 500 * time.Millisecond
	cache := NewCache(5 * time.Millisecond)
	var wg sync.WaitGroup

	started := make(chan struct{})
	wg.Add(1)
	go func() {
		cache.mu.RLock()
		started <- struct{}{}
		defer wg.Done()
		defer cache.mu.RUnlock()
		time.Sleep(waitTime)
	}()
	<-started
	cache.Add("http://example.com", []byte("testdata"))
	if time.Since(initialTime) < waitTime {
		t.Errorf("expected to wait RLock")
		return
	}

	started = make(chan struct{})
	wg.Wait()
	wg.Add(1)
	go func() {
		cache.mu.Lock()
		started <- struct{}{}
		defer wg.Done()
		defer cache.mu.Unlock()
		time.Sleep(waitTime)
	}()
	<-started
	cache.Get("http://example.com")
	if time.Since(initialTime) < 2*waitTime {
		t.Errorf("expected to wait Lock")
		return
	}

	started = make(chan struct{})
	wg.Wait()
	wg.Add(1)
	go func() {
		cache.mu.RLock()
		started <- struct{}{}
		defer wg.Done()
		defer cache.mu.RUnlock()
		time.Sleep(waitTime)
	}()
	<-started
	cache.Get("http://example.com")
	if time.Since(initialTime) > 3*waitTime {
		t.Errorf("expected to not wait RLock")
		return
	}

}
