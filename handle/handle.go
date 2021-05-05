package handle

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/rs/xid"
)

type Handle string

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

func New(prefix string) *Handle {
	handle := Handle(fmt.Sprintf("%s:%s", prefix, xid.New().String()))
	return &handle
}

func NewFor(v interface{}) *Handle {
	return New(prefix(v))
}

func Has(v interface{}) bool {
	rv := reflect.Indirect(reflect.ValueOf(v))
	if rv.Kind() == reflect.Struct && rv.Type().NumField() > 0 && rv.Type().Field(0).Name == "Handle" {
		return true
	}
	return false
}

func Get(v interface{}) *Handle {
	if !Has(v) {
		return nil
	}
	rv := reflect.Indirect(reflect.ValueOf(v))
	h := rv.Field(0).Interface()
	hh, ok := h.(*Handle)
	if !ok {
		return nil
	}
	if hh.Prefix() == "" {
		hh = New(prefix(v))
	}
	return hh
}

// if "" => prefix:
// if id => prefix:id
// if prefix:id => prefix:id
func Set(v interface{}, h string) {
	if !Has(v) {
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
	ptr := reflect.ValueOf(&handle)
	res := reflect.ValueOf(v)
	res.Elem().Field(0).Set(ptr)
}

func prefix(v interface{}) string {
	rv := reflect.Indirect(reflect.ValueOf(v))
	return rv.Type().Field(0).Tag.Get("prefix")
}
