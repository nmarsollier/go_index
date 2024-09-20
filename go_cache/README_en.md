[Versión en Español](README.md)

# A correct way to cache

Caching values sounds something simple to do, but if we do it wrong, it can be a problem.

## A Memoize struct

Lets define some code to store the cached value.

We have a factory function that allows us to cache some value, by retain duration. After that duration value will become nil.

```go
func Memoize(value interface{}, retain time.Duration) *Memo

```

Lets imagine Memo struct (it's in the repo code) will hold the necessary data. We could combine this struct with an interface, but now it's enough.

Memo struct is read only.

If we want to get the cached value, we use Value() function.

```go
func (m *Memo) Value() interface{}
```

This function will return the cached value if it is still valid or nil.

But if we need the cached value, ignoring the validation, we use the Cached() function.

```go
func (m *Memo) Cached() interface{}
```

## A simple usage

Lets assume that we want to have a global cache, for some remote resource that is almost static, but the fetch process takes some time, so we want to create a cache so we don't have to read it each time.

The simplest implementation, **but incorrect** could be:

```go
var cache *memoize.Memo = nil

func wrongCache1(id string) *Profile {
	if cache == nil || cache.Value() == nil {
		value := fetchProfile("123")
		cache = memoize.Memoize(value, 10*time.Minute)
	}

	return cache.Value().(*Profile)
}
```

First we check that the value on cache is valid using Value() function, if not, we fetch data and create the cache.

The we return the cached value.

### Lets see the problems

The previous code could be useful some single thread application, but if we try to run it in a concurrent environment (like rest server engine), we will have many problems...

#### if "if" Race condition

If we close look to the if statement, we can see the first problem :

```go
	if cache == nil || cache.Value() == nil {
```

in a multi thread server the if expression is evaluated in 2 steps, the first one checks for cache == nil, and the seconds gets the Value, but in a race condition the second could raise an error, because cache is not a constant and the value could have become nil in the middle.

For example, a function that invalidates the cache, the most logic would be :

```go
func invalidateCache() {
	cache = nil
}
```

so cache is a global var, because this is the way how it must work.

**Solution**

We do a copy of cache var, in a local value, and we work with that copy, because that is thread safe.

```go
func wrongCache2(id string) *Profile {
	// currCache is thread safe
	currCache := cache

	if currCache == nil || currCache.Value() == nil {
		value := fetchProfile("123")
		currCache = memoize.Memoize(value, 10*time.Minute)
		cache = currCache
	}

	return cache.Value().(*Profile)
}
```

#### return value Race condition

Notice too that the first sample, the return is incorrect, because Value could aso become invalid due to time expiration between check and return.

**Solution**

The correct way to return the value, is using Cahche() function, that does not validates the expiration time.

```go
	return currCache.Cached().(*Profile)
```

#### Data load race condition

Now, fetchProfile is a function that is hard to compute, it takes some time to fetch, and if we have many concurrent process calling this function, those process calling all at the same time will be doing the same highly cost call over the network, the will block responses, and of course could saturate the remote service, making responses more hard to evaluate. [See](https://en.wikipedia.org/wiki/Cache_stampede)

```go
func WrongCache2(id string) *Profile {
	currCache := cache

	if currCache == nil || currCache.Value() == nil {
			value := fetchProfile("123")
```

As we can see in the TestNonConcurrentWrongCache2 we get the wrong desired effect:

```
=== RUN   TestNonConcurrentWrongCache2
Fetching Profile... 0
```

On concurrent calls all the calls are fetching remote values :

```
=== RUN   TestConcurrentWrongCache2
Fetching Profile... 1
Fetching Profile... 5
Fetching Profile... 2
Fetching Profile... 3
Fetching Profile... 4
Fetching Profile... 6
Fetching Profile... 9
Fetching Profile... 0
Fetching Profile... 8
```

**Solution**

We must block the process, so only one process will access to the loading routine at the time, the other process will wait the semaphore.

To do it in a better organized way, i will change the strategy, doing early exists :

```go
// mutex will lock the process
var mutex = &sync.Mutex{}

func wrongCache3(id string) *Profile {
	currCache := cache

	// If the current value is valid, return it
	if currCache != nil && currCache.Value() != nil {
		return currCache.Cached().(*Profile)
	}

	// If the cached value is invalid, we lock the process, so all concurrent calls
	// will wait, and only the first thread will perform the remote call.
	defer(mutex.Unlock())
	mutex.Lock()

	// If other processes waited on the previous lock, cached value now must be
	// valid, so we check that value, and if it is valid, then return it
	currCache = cache
	if currCache != nil && currCache.Value() != nil {
		return currCache.Cached().(*Profile)
	}

	// If this process was not looked before, it's the one that must fetch the data
	value := fetchProfile(id)
	currCache = memoize.Memoize(value, 10*time.Minute)
	cache = currCache

	return currCache.Cached().(*Profile)
}
```

Now we have the desired result, the network call is executed only once.

```
=== RUN   TestConcurrentWrongCache3
Fetching Profile... 2
--- PASS: TestConcurrentWrongCache3 (0.00s)
```

## The final solution

The previous function still has the name wrongCache3, not because it's incorrect, but it's wasting an important opportunity, once we have a cached value, when it expires, we need to reload data, but the value hold in that cache is still valid. With the previous code all concurrent calls will be blocked until we fetch the value, but if the cache has a valid value, a better strategy could be to return that valid value, and in parallel fetch a new data to cache. So concurrent data while loading new value will not be blocked.

What we expect :

- Only perform one network call on concurrent calls
- Only block process when cache is nil
- If cache is expired, but not nil, do not block concurrent calls while loading new data, just return the caced value.

**Solution**

I will explain the final code :

```go
var cache *memoize.Memo = nil
var mutex = &sync.Mutex{}
var loading int32 = 0

func FetchProfile(id string) *Profile {
	// Copy original to create a thread safe var
	currCache := cache
	if currCache != nil && currCache.Value() != nil {
		// Cache is valid, so return Value
		// We return Cached() to avoid expiration race conditions
		return currCache.Cached().(*Profile)
	}

	// If we are here, is because we need to fetch remote data
	// loading is a semaphore, if it has value 1, then we are fetching
	// remote data, a 0 value when we are not in a loading process
	loadData := atomic.CompareAndSwapInt32(&loading, 0, 1)

	// CompareAndSwapInt32 return true if atomically can change from 0 to 1 the
	// value of loading, conceptually I'm assigning a mark to a loading process,
	// so other concurrent threads knows when a fetching is running
	if !loadData && currCache != nil {
		// If we are here, it's because there is a loading process running
		// and the currCache is not nil, so the value is valid
		return currCache.Cached().(*Profile)
	}

	// Until now, we where not blocking concurrent process, now we need to
	// ensure that only one process will be loading data at the time
	defer mutex.Unlock()
	mutex.Lock()

	// If we where blocked by the previous Lock, when we wake up, the cache
	// could hold the value to return, we need to ensure use it to prevent
	// incorrect fetching
	currCache = cache
	if loading == 0 && currCache != nil && currCache.Value() != nil {
		// We are here if we have been Looked previously
		return currCache.Cached().(*Profile)
	}

	// Set that we are not loading anymore
	// we cannot do it before because there are returns that would set it
	// 0 incorrectly
	defer func() { loading = 0 }()

	// Fetch the data
	value := fetchProfile(id)

	// Update cache and return
	currCache = memoize.Memoize(value, 10*time.Minute)
	cache = currCache

	return currCache.Cached().(*Profile)
}
```

**Tests that validates the process**

The next test is weird, just to visually demonstrates the correct working :

```go
func fineFetchProfile(t *testing.T) {
	invalidateCache()

	// 2 steps, first we call 10 concurrent calls with cache = nil
	// first call should fetch, the rest should lock and wait
	var waitGroup sync.WaitGroup
	waitGroup.Add(10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			defer waitGroup.Done()
			p := fineFetchProfile(strconv.Itoa(i))
			t.Logf("Result Step 1 = %s = %s \n", strconv.Itoa(i), p.Name)
		}(i)
	}
	waitGroup.Wait()

	// Previous calls should work until cache expire, i will simulate
	// a cache expiration, setting a value that expires in 1 second
	cache = memoize.Memoize(fetchProfile("Expired"), 1*time.Second)
	time.Sleep(2 * time.Second)

	// Next 10 concurrent calls, are done with an expired cache, but valid value
	// we spect that the first call fetch the data, but the others does not
	// lock, but return cached value
	waitGroup.Add(10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			defer waitGroup.Done()
			p := fineFetchProfile(strconv.Itoa(i))
			t.Logf("Result Step 2 = %s = %s \n", strconv.Itoa(i), p.Name)
		}(i)
	}
	waitGroup.Wait()

	p := fineFetchProfile("Final")
	t.Logf("Value after changes = %s \n", p.Name)

}
```

We can see in the results :

```
=== RUN   TestConcurrentFetchProfile
Fetching Profile... 1        // fetch call 1

// There is no cahce, first call fetchs and is the first to return
TestConcurrentFetchProfile: service_test.go:50: Result Step 1 = 1 = Profile # 1
TestConcurrentFetchProfile: service_test.go:50: Result Step 1 = 5 = Profile # 1
TestConcurrentFetchProfile: service_test.go:50: Result Step 1 = 2 = Profile # 1
TestConcurrentFetchProfile: service_test.go:50: Result Step 1 = 9 = Profile # 1
TestConcurrentFetchProfile: service_test.go:50: Result Step 1 = 0 = Profile # 1
TestConcurrentFetchProfile: service_test.go:50: Result Step 1 = 3 = Profile # 1
TestConcurrentFetchProfile: service_test.go:50: Result Step 1 = 4 = Profile # 1
TestConcurrentFetchProfile: service_test.go:50: Result Step 1 = 8 = Profile # 1
TestConcurrentFetchProfile: service_test.go:50: Result Step 1 = 6 = Profile # 1
TestConcurrentFetchProfile: service_test.go:50: Result Step 1 = 7 = Profile # 1

// When we have an expired cache...
// first call will fetch data
Fetching Profile... 1

// the concurrent calls will receive the caced value
TestConcurrentFetchProfile: service_test.go:63: Result Step 2 = 2 = Profile # Expired
TestConcurrentFetchProfile: service_test.go:63: Result Step 2 = 6 = Profile # Expired
TestConcurrentFetchProfile: service_test.go:63: Result Step 2 = 4 = Profile # Expired
TestConcurrentFetchProfile: service_test.go:63: Result Step 2 = 0 = Profile # Expired
TestConcurrentFetchProfile: service_test.go:63: Result Step 2 = 7 = Profile # Expired
TestConcurrentFetchProfile: service_test.go:63: Result Step 2 = 9 = Profile # Expired
TestConcurrentFetchProfile: service_test.go:63: Result Step 2 = 3 = Profile # Expired
TestConcurrentFetchProfile: service_test.go:63: Result Step 2 = 5 = Profile # Expired
TestConcurrentFetchProfile: service_test.go:63: Result Step 2 = 8 = Profile # Expired
TestConcurrentFetchProfile: service_test.go:63: Result Step 2 = 1 = Profile # 1
TestConcurrentFetchProfile: service_test.go:69: Value after changes = Profile # 0

// After process 1 fetchs data will return at the end
--- PASS: TestConcurrentFetchProfile
```

## Doing a library to generalize the logic

---

NOTE

We are missing another important oportunity in the before routine, and the question is why on the second step I wait to return the new value for profile 1, why I can't return the cached value for that request too, and fetch data in a goroutine ?

This final library solves that problem, but I will let you find how.

---

SafeMemoize is an struct defined in file safe_memorize.go that allow us to generalize the cache logic.

It has only one method, besides the constructor :

```go
// Value get cached value, fetching data if needed
func (m *SafeMemoize) Value(
	fetchFunc func() *Memo,
) interface{} {
	...
}
```

The only piece of code that we need to privide is the updater function, that most if the times it will be a closure like this :

```go
var profileMemoize = memoize.NewSafeMemoize()

// FetchProfile fetch the current profile
func FetchProfile(id string) *Profile {
	return profileMemoize.Value(
		func() *memoize.Memo {
			return memoize.Memoize(fetchProfile(id), 10*time.Minute)
		},
	).(*Profile)
}
```

And that's it, we have a cache library to catch any value from remotes, and we can use it easy.

## Note

This is a series of notes about advanced Go patterns, with a really simple implementation.

[Content Table](https://github.com/nmarsollier/go_index/blob/main/README_en.md)
