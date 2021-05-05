package handle

import (
	"testing"
)

type FakeResource struct {
	*Handle `prefix:"res"`

	Foo string
	Bar int
}

func TestHandle(t *testing.T) {
	v := FakeResource{}

	if !Has(v) {
		t.Fatal("resource expected to have handle")
	}
	h := Get(v)
	if h.Prefix() != "res" {
		t.Fatal("resource handle expected to have prefix 'res'")
	}
	if h.ID() == "" {
		t.Fatal("resource handle expected to have non-empty id")
	}

	nh := New("res")
	if nh.Prefix() != "res" {
		t.Fatal("new handle expected to have prefix 'res'")
	}
	if nh.ID() == "" {
		t.Fatal("new handle expected to have non-empty id")
	}

	Set(&v, "")
	if Get(v).Prefix() != "res" {
		t.Fatal("set handle expected to have prefix 'res'")
	}
	if Get(v).ID() != "" {
		t.Fatal("set handle expected to have empty id")
	}

	Set(&v, "123")
	if Get(v).Prefix() != "res" {
		t.Fatal("set handle expected to have prefix 'res'")
	}
	if Get(v).ID() != "123" {
		t.Fatal("set handle expected to have id '123'")
	}

	Set(&v, "res:123")
	if Get(v).Prefix() != "res" {
		t.Fatal("set handle expected to have prefix 'res'")
	}
	if Get(v).ID() != "123" {
		t.Fatal("set handle expected to have id '123'")
	}

}
