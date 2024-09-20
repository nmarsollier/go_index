package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSayHelo(t *testing.T) {
	defaultHelloFunc := daoHelloFunc

	// Cuando testeamos la reescribimos con el
	// mock que queramos
	daoHelloFunc = func() string {
		return "Hello"
	}

	assert.Equal(t, "Hello Pepe", SayHello("Pepe"))
	daoHelloFunc = defaultHelloFunc
}
