package resource

import (
	"testing"
)

type FakeResource struct {
	*Handle `prefix:"res"`

	Foo string
	Bar int
}

func TestResource(t *testing.T) {
	Register(FakeResource{})
	v := New("res")

	if !HasHandle(v) {
		t.Fatal("resource expected to have handle")
	}

	h := GetHandle(v)
	if h.Prefix() != "res" {
		t.Fatal("resource handle expected to have prefix 'res'")
	}
	if h.ID() == "" {
		t.Fatal("resource handle expected to have non-empty id")
	}

	nh := NewHandle("res")
	if nh.Prefix() != "res" {
		t.Fatal("new handle expected to have prefix 'res'")
	}
	if nh.ID() == "" {
		t.Fatal("new handle expected to have non-empty id")
	}

	SetHandle(v, "")
	if GetHandle(v).Prefix() != "res" {
		t.Fatal("set handle expected to have prefix 'res'")
	}
	if GetHandle(v).ID() != "" {
		t.Fatal("set handle expected to have empty id")
	}

	SetHandle(v, "123")
	if GetHandle(v).Prefix() != "res" {
		t.Fatal("set handle expected to have prefix 'res'")
	}
	if GetHandle(v).ID() != "123" {
		t.Fatal("set handle expected to have id '123'")
	}

	SetHandle(v, "res:123")
	if GetHandle(v).Prefix() != "res" {
		t.Fatal("set handle expected to have prefix 'res'")
	}
	if GetHandle(v).ID() != "123" {
		t.Fatal("set handle expected to have id '123'")
	}

}
