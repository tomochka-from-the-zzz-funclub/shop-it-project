package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

var errNotFound = errors.New("constructor for such type not found")
var errWrongTypeForConstructor = errors.New("constructor record have wrong type")

type UserService struct {
	NotEmptyStruct bool
}
type MessageService struct {
	NotEmptyStruct bool
}

type Container struct {
	constructors map[string]interface{}
}

func NewContainer() *Container {
	return &Container{constructors: make(map[string]interface{})}
}

func (c *Container) RegisterType(name string, constructor interface{}) {
	c.constructors[name] = constructor
}

func (c *Container) Resolve(name string) (interface{}, error) {
	constructor, found := c.constructors[name]
	if !found {
		return nil, errNotFound
	}
	f, ok := constructor.(func() interface{})
	if !ok {
		return nil, errWrongTypeForConstructor
	}
	return f(), nil
}

func TestDIContainer(t *testing.T) {
	container := NewContainer()
	container.RegisterType("UserService", func() interface{} {
		return &UserService{}
	})
	container.RegisterType("MessageService", func() interface{} {
		return &MessageService{}
	})

	userService1, err := container.Resolve("UserService")
	assert.NoError(t, err)
	userService2, err := container.Resolve("UserService")
	assert.NoError(t, err)

	u1 := userService1.(*UserService)
	u2 := userService2.(*UserService)
	assert.False(t, u1 == u2)

	messageService, err := container.Resolve("MessageService")
	assert.NoError(t, err)
	assert.NotNil(t, messageService)

	paymentService, err := container.Resolve("PaymentService")
	assert.Error(t, err)
	assert.Nil(t, paymentService)
}
