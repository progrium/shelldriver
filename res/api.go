package res

import (
	"fmt"

	"github.com/progrium/macbridge/handle"
)

func (m *Manager) Sync(v interface{}) error {
	if !handle.Has(v) {
		return fmt.Errorf("not a resource")
	}
	h := handle.Get(v)
	if h == nil {
		h = handle.NewFor(v)
	}
	var rh string
	_, err := m.Peer.Call("Sync", []interface{}{*h, v}, &rh)
	handle.Set(v, rh)
	return err
}

func (m *Manager) Release(v interface{}) error {
	if !handle.Has(v) {
		return fmt.Errorf("not a resource")
	}
	h := handle.Get(v)
	if h == nil {
		return fmt.Errorf("unable to release an uninitialized resource")
	}
	_, err := m.Peer.Call("Release", *h, nil)
	if err == nil {
		handle.Set(v, "")
	}
	return err
}
