package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSayHelo(t *testing.T) {
	defaultHelloFunc := daoGetHello

	// Cuando testeamos la reescribimos con el
	// mock que queramos
	daoGetHello = func() string {
		return "Hello"
	}

	assert.Equal(t, "Hello Pepe", SayHello("Pepe"))
	daoGetHello = defaultHelloFunc
}
