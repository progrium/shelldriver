package resource

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/rs/xid"
)

var registeredTypes map[string]reflect.Type

func init() {
	registeredTypes = make(map[string]reflect.Type)
}

func TypePrefix(v interface{}) string {
	rt := reflect.Indirect(reflect.ValueOf(v)).Type()
	for p, t := range registeredTypes {
		if t == rt {
			return p
		}
	}
	log.Panicf("type '%v' not registered resource", rt)
	return ""
}

func Register(v interface{}) {
	rv := reflect.Indirect(reflect.ValueOf(v))
	prefix := rv.Type().Field(0).Tag.Get("prefix")
	registeredTypes[prefix] = rv.Type()
}

func New(prefix string) interface{} {
	t, ok := registeredTypes[prefix]
	if !ok {
		log.Panicf("resource type not registered: %s", prefix)
	}
	h := NewHandle(prefix)
	r := reflect.New(t).Interface()
	SetHandle(r, h.Handle())
	return r
}

type Handle string

func NewHandle(prefix string) *Handle {
	handle := Handle(fmt.Sprintf("%s:%s", prefix, xid.New().String()))
	return &handle
}

func HasHandle(v interface{}) bool {
	rv := reflect.Indirect(reflect.ValueOf(v))
	if rv.Kind() == reflect.Struct && rv.Type().NumField() > 0 && rv.Type().Field(0).Name == "Handle" {
		return true
	}
	return false
}

func GetHandle(v interface{}) *Handle {
	if !HasHandle(v) {
		return nil
	}
	rv := reflect.Indirect(reflect.ValueOf(v))
	h := rv.Field(0).Interface()
	hh, ok := h.(*Handle)
	if !ok {
		return nil
	}
	if hh.Prefix() == "" {
		hh = NewHandle(TypePrefix(v))
	}
	return hh
}

// if "" => prefix:
// if id => prefix:id
// if prefix:id => prefix:id
func SetHandle(v interface{}, h string) {
	if !HasHandle(v) {
		return
	}
	if !strings.Contains(h, ":") {
		if h == "" {
			h = fmt.Sprintf("%s:", TypePrefix(v))
		} else {
			h = fmt.Sprintf("%s:%s", TypePrefix(v), h)
		}
	}
	handle := Handle(h)
	ptr := reflect.ValueOf(&handle)
	res := reflect.ValueOf(v)
	res.Elem().Field(0).Set(ptr)
}

func (h *Handle) Prefix() string {
	if h == nil {
		return ""
	}
	parts := strings.Split(string(*h), ":")
	return parts[0]
}

func (h *Handle) ID() string {
	if h == nil {
		return ""
	}
	parts := strings.Split(string(*h), ":")
	if len(parts) > 1 {
		return parts[1]
	}
	return ""
}

func (h *Handle) Handle() string {
	if h == nil {
		return ""
	}
	return string(*h)
}
