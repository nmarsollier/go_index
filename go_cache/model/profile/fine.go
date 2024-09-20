package profile

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/nmarsollier/go_cache/utils/memoize"
)

var cache *memoize.Memo = nil
var mutex = &sync.Mutex{}
var loading int32 = 0

// FetchProfile fetch the current profile
func fineFetchProfile(id string) *Profile {
	currCache := cache
	if currCache != nil && currCache.Value() != nil {
		return currCache.Cached().(*Profile)
	}

	loadData := atomic.CompareAndSwapInt32(&loading, 0, 1)
	if !loadData && currCache != nil {
		return currCache.Cached().(*Profile)
	}

	defer mutex.Unlock()
	mutex.Lock()
	currCache = cache
	if loading == 0 && currCache != nil && currCache.Value() != nil {
		return currCache.Cached().(*Profile)
	}

	defer func() {
		loading = 0
	}()

	value := fetchProfile(id)
	currCache = memoize.Memoize(value, 10*time.Minute)
	cache = currCache

	return currCache.Cached().(*Profile)
}

func invalidateCache() {
	cache = nil
}
