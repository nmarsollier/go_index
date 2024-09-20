package memoize

import "time"

type Memo struct {
	created time.Time
	retain  time.Duration
	expire  time.Time
	value   interface{}
}

// Memoize a value for the given time, retain = 0 means forever
func Memoize(value interface{}, retain time.Duration) *Memo {
	return &Memo{
		retain:  retain,
		created: time.Now(),
		expire:  time.Now().Add(retain),
		value:   value,
	}
}

func timestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// Value is the cached value or nil if no longer valid
func (m *Memo) Value() interface{} {
	if m.retain == 0 {
		return m.value
	}

	if time.Now().Before(m.expire) {
		return m.value
	}

	return nil
}

// Cached value
func (m *Memo) Cached() interface{} {
	return m.value
}
