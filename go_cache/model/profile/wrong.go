package profile

import (
	"time"

	"github.com/nmarsollier/go_cache/utils/memoize"
)

func wrongCache1(id string) *Profile {
	if cache == nil || cache.Value() == nil {
		value := fetchProfile(id)
		cache = memoize.Memoize(value, 10*time.Minute)
	}

	return cache.Value().(*Profile)
}

func wrongCache2(id string) *Profile {
	currCache := cache

	if currCache == nil || currCache.Value() == nil {
		value := fetchProfile(id)
		currCache = memoize.Memoize(value, 10*time.Minute)
		cache = currCache
	}

	return currCache.Cached().(*Profile)
}

func wrongCache3(id string) *Profile {
	currCache := cache
	if currCache != nil && currCache.Value() != nil {
		return currCache.Cached().(*Profile)
	}

	defer mutex.Unlock()
	mutex.Lock()
	currCache = cache
	if currCache != nil && currCache.Value() != nil {
		return currCache.Cached().(*Profile)
	}

	value := fetchProfile(id)
	currCache = memoize.Memoize(value, 10*time.Minute)
	cache = currCache

	return currCache.Cached().(*Profile)
}
