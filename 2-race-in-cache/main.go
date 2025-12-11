//////////////////////////////////////////////////////////////////////
//
// Given is some code to cache key-value pairs from a database into
// the main memory (to reduce access time). Note that golang's map are
// not entirely thread safe. Multiple readers are fine, but multiple
// writers are not. Change the code to make this thread safe.
//

package main

import (
	"container/list"
	"sync" 
	"testing"
)

// CacheSize determines how big the cache can grow
const CacheSize = 100

// KeyStoreCacheLoader is an interface for the KeyStoreCache
type KeyStoreCacheLoader interface {
	// Load implements a function where the cache should gets it's content from
	Load(string) string
}

type page struct {
	Key   string
	Value string
}

// KeyStoreCache is a LRU cache for string key-value pairs
type KeyStoreCache struct {
	mu    sync.RWMutex
	cache map[string]*list.Element
	pages list.List
	load  func(string) string
}

// New creates a new KeyStoreCache
func New(load KeyStoreCacheLoader) *KeyStoreCache {
	return &KeyStoreCache{
		load:  load.Load,
		cache: make(map[string]*list.Element),
	}
}

// Get gets the key from cache, loads it from the source if needed
func (k *KeyStoreCache) Get(key string) string {
	// initial read lock
	k.mu.RLock()
	if e, ok := k.cache[key]; ok {
		// if cache hit, unlock read lock then upgrade to write lock for 'MoveToFront' func
		// since we're gonna move an element. list.List methods are not thread safe
		k.mu.RUnlock()
		
		k.mu.Lock()
		k.pages.MoveToFront(e)
		k.mu.Unlock()

		return e.Value.(page).Value
	}
	k.mu.RUnlock()

	// if cache miss then we're writing to we're gonna write lock
	k.mu.Lock()
	defer k.mu.Unlock()

	// key might have been loaded by another goroutine so we check again while this goroutine is waiting for the lock.
	if e, ok := k.cache[key]; ok {	// Key was loaded by another goroutine, just move it to front and return.
		k.pages.MoveToFront(e)
		return e.Value.(page).Value
	}

	// Miss - load from database and save it in cache
	// this is the expensive operation we are protecting against with the double check.
	p := page{key, k.load(key)}
	
	// if cache is full remove the least used item
	if len(k.cache) >= CacheSize {
		end := k.pages.Back()
		// remove from map
		delete(k.cache, end.Value.(page).Key)
		// remove from list
		k.pages.Remove(end)
	}

	k.pages.PushFront(p)
	k.cache[key] = k.pages.Front()
	
	return p.Value
}

// Loader implements KeyStoreLoader
type Loader struct {
	DB *MockDB
}

// Load gets the data from the database
func (l *Loader) Load(key string) string {
	val, err := l.DB.Get(key)
	if err != nil {
		panic(err)
	}

	return val
}

func run(t *testing.T) (*KeyStoreCache, *MockDB) {
	loader := Loader{
		DB: GetMockDB(),
	}
	cache := New(&loader)

	RunMockServer(cache, t)

	return cache, loader.DB
}

func main() {
	run(nil)
}
