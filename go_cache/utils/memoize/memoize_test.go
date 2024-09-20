package memoize

import (
	"testing"
	"time"

	"gopkg.in/go-playground/assert.v1"
)

func TestMemoize1Sec(t *testing.T) {
	memo1Sec := Memoize("hello", 1000)

	assert.Equal(t, memo1Sec.Value(), "hello")
	assert.Equal(t, memo1Sec.Cached(), "hello")

	time.Sleep(1 * time.Second)

	assert.Equal(t, memo1Sec.Value(), nil)
	assert.Equal(t, memo1Sec.Cached(), "hello")
}
