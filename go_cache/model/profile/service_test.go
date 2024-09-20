package profile

import (
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/nmarsollier/go_cache/utils/memoize"
)

func TestNonConcurrentWrongCache2(t *testing.T) {
	invalidateCache()
	for i := 0; i < 10; i++ {
		wrongCache2(strconv.Itoa(i))
	}
}

func TestConcurrentWrongCache2(t *testing.T) {
	invalidateCache()
	for i := 0; i < 10; i++ {
		go wrongCache2(strconv.Itoa(i))
	}

	for i := 0; i < 10; i++ {
		go wrongCache2(strconv.Itoa(i))
	}
}

func TestConcurrentWrongCache3(t *testing.T) {
	invalidateCache()
	for i := 0; i < 10; i++ {
		go wrongCache3(strconv.Itoa(i))
	}

	for i := 0; i < 10; i++ {
		go wrongCache3(strconv.Itoa(i))
	}
}

func TestConcurrentFetchProfile(t *testing.T) {
	invalidateCache()

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

	cache = memoize.Memoize(fetchProfile("Expired"), 1*time.Second)
	time.Sleep(2 * time.Second)

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

func TestSafeFetchProfile(t *testing.T) {
	invalidateTSCache()

	var waitGroup sync.WaitGroup
	waitGroup.Add(10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			defer waitGroup.Done()
			p := FetchProfile(strconv.Itoa(i))
			t.Logf("Result Step 1 = %s = %s \n", strconv.Itoa(i), p.Name)
		}(i)
	}
	waitGroup.Wait()

	profileMemoize.ReplaceMockCache(memoize.Memoize(fetchProfile("Expired"), 1*time.Second))
	time.Sleep(2 * time.Second)

	waitGroup.Add(10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			defer waitGroup.Done()
			p := FetchProfile(strconv.Itoa(i))
			t.Logf("Result Step 2 = %s = %s \n", strconv.Itoa(i), p.Name)
		}(i)
	}
	waitGroup.Wait()

	// Lets wait until fetch goroutine ends
	time.Sleep(2 * time.Second)

	p := FetchProfile("Final")
	t.Logf("Value after changes = %s \n", p.Name)
}
