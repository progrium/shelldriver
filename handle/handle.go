package handle

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/rs/xid"
)

const Invalid = Handle("")

type Resourcer interface {
	Resource() (*Handle, interface{})
}

type Handle string

func (h Handle) Type() string {
	parts := strings.Split(string(h), ":")
	return parts[0]
}

func (h Handle) ID() string {
	parts := strings.Split(string(h), ":")
	if len(parts) > 1 {
		return parts[1]
	}
	return ""
}

func (h Handle) Unset() bool {
	return h.ID() == ""
}

func (h Handle) Handle() string {
	return string(h)
}

func New(typ, id string) Handle {
	return Handle(fmt.Sprintf("%s:%s", typ, id))
}

func NewFor(v interface{}) Handle {
	return New(prefix(v), xid.New().String())
}

func Has(v interface{}) bool {
	if _, ok := v.(Resourcer); ok {
		return true
	}
	rv := reflect.Indirect(reflect.ValueOf(v))
	if rv.Kind() == reflect.Struct && rv.Type().NumField() > 0 &&
		rv.Type().Field(0).Name == "Handle" && rv.Type().Field(0).Type.Name() == "Handle" {
		return true
	}
	return false
}

func Get(v interface{}) (handle Handle) {
	if r, ok := v.(Resourcer); ok {
		h, _ := r.Resource()
		handle = *h
	} else {
		rv := reflect.Indirect(reflect.ValueOf(v))
		h := rv.Field(0).Interface()
		var ok bool
		handle, ok = h.(Handle)
		if !ok {
			return Invalid
		}
	}
	if handle.Type() == "" {
		handle = New(prefix(v), "")
	}
	return handle
}

// if "" => type:
// if id => type:id
// if type:id => type:id
func Set(v interface{}, h string) {
	if !strings.Contains(h, ":") {
		if h == "" {
			h = fmt.Sprintf("%s:", prefix(v))
		} else {
			h = fmt.Sprintf("%s:%s", prefix(v), h)
		}
	}
	sethandle := Handle(h)
	if r, ok := v.(Resourcer); ok {
		h, _ := r.Resource()
		*h = sethandle
	} else {
		ptr := reflect.ValueOf(sethandle)
		res := reflect.ValueOf(v)
		res.Elem().Field(0).Set(ptr)
	}
}

func prefix(v interface{}) string {
	rv := reflect.Indirect(reflect.ValueOf(v))
	prefixTag := rv.Type().Field(0).Tag.Get("type")
	if prefixTag == "" {
		return rv.Type().Name()
	}
	return prefixTag
}
