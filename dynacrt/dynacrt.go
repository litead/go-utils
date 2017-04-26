package dynacrt

import (
	"fmt"
)

type CreateFunc func(arg interface{}) (interface{}, error)

var registry = make(map[string]CreateFunc)

func Register(typeName string, createFn CreateFunc) {
	if registry[typeName] != nil {
		panic(fmt.Errorf("type '%s' already registered", typeName))
	}
	registry[typeName] = createFn
}

func Create(typeName string, arg interface{}) (interface{}, error) {
	fn := registry[typeName]
	if fn == nil {
		return nil, fmt.Errorf("type '%s' not registered", typeName)
	}
	return fn(arg)
}
