package handle

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/rs/xid"
)

const Invalid = Handle("")

type Resourcer interface {
	Resource() interface{}
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

func (h Handle) IsZero() bool {
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
		v = r.Resource()
	}
	rv := reflect.Indirect(reflect.ValueOf(v))
	h := rv.Field(0).Interface()
	var ok bool
	handle, ok = h.(Handle)
	if !ok {
		return Invalid
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
	if r, ok := v.(Resourcer); ok {
		Set(r.Resource(), h)
		return
	}
	if !strings.Contains(h, ":") {
		if h == "" {
			h = fmt.Sprintf("%s:", prefix(v))
		} else {
			h = fmt.Sprintf("%s:%s", prefix(v), h)
		}
	}
	handle := Handle(h)
	rhandle := reflect.ValueOf(handle)
	res := reflect.ValueOf(v)
	res.Elem().Field(0).Set(rhandle)
}

func prefix(v interface{}) string {
	rv := reflect.Indirect(reflect.ValueOf(v))
	prefixTag := rv.Type().Field(0).Tag.Get("type")
	if prefixTag == "" {
		return rv.Type().Name()
	}
	return prefixTag
}
