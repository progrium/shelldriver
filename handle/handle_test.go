package handle

import (
	"testing"
)

type Nonresource struct {
	Handle int
}

type Fake struct {
	Handle

	Foo string
}

type NamedFake struct {
	Handle `type:"fak"`

	Foo string
	Bar int
}

type EmbeddedFake struct {
	Fake
}

func (ef *EmbeddedFake) Resource() interface{} {
	return &ef.Fake
}

func TestEmbedded(t *testing.T) {
	v := &EmbeddedFake{}

	if !Has(v) {
		t.Fatal("resource expected to have handle")
	}

	h := Get(v)
	if !h.IsZero() {
		t.Fatal("resource handle expected to be zero")
	}

	Set(v, "123")
	h = Get(v)
	if h.ID() != "123" {
		t.Fatal("resource expected to have id 123")
	}

}

func TestNonresource(t *testing.T) {
	v := &Nonresource{}

	if Has(v) {
		t.Fatal("resource expected to not have handle")
	}

	h := Get(v)
	if h != Invalid {
		t.Fatal("resource handle expected to be invalid")
	}
	if h.ID() != "" {
		t.Fatal("resource handle expected to have empty id")
	}
}

func TestUnnamedHandle(t *testing.T) {
	v := &Fake{}

	if !Has(v) {
		t.Fatal("resource expected to have handle")
	}

	h := Get(v)
	if !h.IsZero() {
		t.Fatal("resource handle expected to be zero")
	}

	Set(v, "123")
	h = Get(v)
	if h.Type() != "Fake" {
		t.Fatal("resource handle expected to have type 'Fake'")
	}
	id1 := h.ID()
	id2 := Get(v).ID()
	if id1 != id2 {
		t.Fatal("resource handle expected to be consistent: ", id1, id2)
	}

}

func TestHandle(t *testing.T) {
	v := &NamedFake{}

	if !Has(v) {
		t.Fatal("resource expected to have handle")
	}
	h := Get(v)
	if !h.IsZero() {
		t.Fatal("resource handle expected to be zero")
	}

	nh := New("res", "")
	if nh.Type() != "res" {
		t.Fatal("new handle expected to have type 'res'")
	}
	if nh.ID() != "" {
		t.Fatal("new handle expected to have empty id")
	}

	Set(v, "")
	if Get(v).Type() != "fak" {
		t.Fatal("set handle expected to have type 'fak'")
	}
	if Get(v).ID() != "" {
		t.Fatal("set handle expected to have empty id")
	}

	Set(v, "123")
	if Get(v).Type() != "fak" {
		t.Fatal("set handle expected to have type 'fak'")
	}
	if Get(v).ID() != "123" {
		t.Fatal("set handle expected to have id '123'")
	}

	Set(v, "res:123")
	if Get(v).Type() != "res" {
		t.Fatal("set handle expected to have type 'res'")
	}
	if Get(v).ID() != "123" {
		t.Fatal("set handle expected to have id '123'")
	}

}
