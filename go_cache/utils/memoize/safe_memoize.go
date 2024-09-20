package memoize

import (
	"sync"
	"sync/atomic"
)

type SafeMemoize struct {
	cache   *Memo
	mutex   *sync.Mutex
	loading int32
}

// InvalidateCache invalidates the cache
func (m *SafeMemoize) InvalidateCache() {
	m.cache = nil
}

// ReplaceMockCache just to mock test, this should be removed
func (m *SafeMemoize) ReplaceMockCache(newCache *Memo) {
	m.cache = newCache
}

// Value get cached value, fetching data if needed
func (m *SafeMemoize) Value(
	fetchFunc func() *Memo,
) interface{} {
	currCache := m.cache
	if currCache != nil && currCache.Value() != nil {
		return currCache.Cached()
	}

	loadData := atomic.CompareAndSwapInt32(&m.loading, 0, 1)
	if currCache == nil {
		// No data cached, just lock the concurrent calls
		currCache = m.fetchData(fetchFunc)
	} else if loadData {
		// There is a valid cache, load in goroutine
		go m.fetchData(fetchFunc)
	}

	// The only possibility of nil, is on first cache load error
	// process will retry but this one has failed
	if currCache == nil {
		return nil
	}

	return currCache.Cached()
}

func (m *SafeMemoize) fetchData(
	fetchFunc func() *Memo,
) *Memo {
	defer m.mutex.Unlock()
	m.mutex.Lock()
	currCache := m.cache
	if m.loading == 0 && currCache != nil && currCache.Value() != nil {
		return currCache
	}

	defer func() { m.loading = 0 }()

	currCache = fetchFunc()
	if currCache != nil {
		// nil means an error, to be resilent will not update chache
		// dev need to work on a good error handler
		m.cache = currCache
	}
	return currCache
}

// NewSafeMemoize creates new thread safe memoization
func NewSafeMemoize() *SafeMemoize {
	return &SafeMemoize{
		cache:   nil,
		mutex:   &sync.Mutex{},
		loading: 0,
	}
}
