package bridge

import (
	"sync"

	"github.com/mitchellh/mapstructure"
	"github.com/progrium/macbridge/handle"
)

type Resource interface {
	// Resource returns a pointer to the handled struct
	Resource() interface{}

	// Apply attempts to apply the handled struct properties
	// to the actual shell managed resource
	Apply() error

	// Discard attempts to close and release the
	// shell managed resource
	Discard() error
}

type Bridge struct {
	resources map[handle.Handle]Resource
	sync.Mutex
}

func New() *Bridge {
	return &Bridge{
		resources: make(map[handle.Handle]Resource),
	}
}

// Types returns the names of supported resource types
func (b *Bridge) Types() []string {
	var t []string
	for k := range types {
		t = append(t, k)
	}
	return t
}

// Discard closes and releases the shell managed resource by the given handle
func (b *Bridge) Discard(h handle.Handle) error {
	b.Lock()
	defer b.Unlock()
	r, ok := b.resources[h]
	if !ok {
		return nil
	}
	defer delete(b.resources, h)
	return Dispatch(r.Discard)
}

// Apply either creates a shell managed resource of the given handle type using
// the properties of the map or structure provided, or it applies the properties
// of the given map or structure to the managed resource by the given handle
func (b *Bridge) Apply(h handle.Handle, res interface{}) (Resource, error) {
	b.Lock()
	r, ok := b.resources[h]
	if !ok {
		h, r = new_(h.Type())
		b.resources[h] = r
	}
	if err := mapstructure.Decode(res, r); err != nil {
		b.Unlock()
		return nil, err
	}
	handle.Set(r, h.Handle())
	b.Unlock()
	return r, Dispatch(r.Apply)
}

// Sync takes a pointer to a handled struct, runs it through Apply with
// its own handle, then applies the resulting resource properties back
// to the handled struct pointer given
func (b *Bridge) Sync(res interface{}) error {
	newres, err := b.Apply(handle.Get(res), res)
	if err != nil {
		return err
	}
	return mapstructure.Decode(newres, res)
}

// Close discards all known resources
func (b *Bridge) Close() (err error) {
	b.Lock()
	defer b.Unlock()
	for h := range b.resources {
		b.Unlock()
		err = b.Discard(h)
		b.Lock()
		if err != nil {
			return
		}
	}
	return
}
