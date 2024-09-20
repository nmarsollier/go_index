package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type daoMock struct {
}

func (d daoMock) Hello() string {
	return "Hello"
}

func TestSayHelo(t *testing.T) {
	// Mockeamos
	mockedDao := new(daoMock)

	s := HelloService{
		mockedDao,
	}
	assert.Equal(t, "Hello", s.SayHello())
}
