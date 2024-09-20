package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/**
En nuestro test tendremos que mockear el dao para pasarlo como DI
*/
type daoMock struct {
}

func (d daoMock) Hello() string {
	return "Hello"
}

func TestSayHelo(t *testing.T) {
	s := NewService(new(daoMock))

	assert.Equal(t, "Hello", s.SayHello())
}
