package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSayHelo(t *testing.T) {
	// Cuando testeamos la reescribimos con el
	// mock que queramos
	sayHelloFunc = func() string {
		return "Hello"
	}

	assert.Equal(t, "Hello", SayHello())
}
