package bridge

import (
	"reflect"

	"github.com/progrium/shelldriver/handle"
)

var types map[string]reflect.Type

// register makes a given resource type available to the bridge
func register(v Resource) {
	if types == nil {
		types = make(map[string]reflect.Type)
	}
	types[handle.Get(v).Type()] = reflect.TypeOf(v).Elem()
}

// new_ returns a new instance of a resource type along with its handle
func new_(typ string) (handle.Handle, Resource) {
	t, ok := types[typ]
	if !ok {
		panic("bridge type not available: " + typ)
	}
	r := reflect.New(t).Interface()
	handle.Set(r, handle.NewFor(r).Handle())
	return handle.Get(r), r.(Resource)
}
