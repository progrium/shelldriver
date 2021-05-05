package bridge

import (
	"context"
	"fmt"
	"io"
	"os"
	"reflect"
	"sync"

	"github.com/manifold/qtalk/golang/mux"
	"github.com/manifold/qtalk/golang/rpc"
	"github.com/mitchellh/mapstructure"
	"github.com/progrium/macbridge/handle"
	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
)

type Bridge struct {
	Resources []interface{}
	Types     map[string]reflect.Type

	objects  map[handle.Handle]objc.Object
	released map[handle.Handle]bool

	sync.Mutex
}

func New() *Bridge {
	b := &Bridge{
		Types:    make(map[string]reflect.Type),
		objects:  make(map[handle.Handle]objc.Object),
		released: make(map[handle.Handle]bool),
	}
	b.Register(Window{})
	b.Register(Indicator{})
	b.Register(Menu{})
	return b
}

func (b *Bridge) Run() {
	app := cocoa.NSApp_WithDidLaunch(func(n objc.Object) {
		session := mux.NewSession(context.Background(), struct {
			io.ReadCloser
			io.Writer
		}{os.Stdin, os.Stdout})
		peer := rpc.NewPeer(session, rpc.JSONCodec{})
		peer.Bind("Sync", b.Sync)
		peer.Bind("Release", b.Release)
		go peer.Respond()
	})
	app.SetActivationPolicy(cocoa.NSApplicationActivationPolicyRegular)
	app.ActivateIgnoringOtherApps(true)
	app.Run()
}

func (b *Bridge) Log(v ...interface{}) {
	fmt.Fprintln(os.Stderr, v...)
}

func (b *Bridge) Register(v interface{}) {
	rv := reflect.Indirect(reflect.ValueOf(v))
	prefix := rv.Type().Field(0).Tag.Get("prefix")
	b.Types[prefix] = rv.Type()
}

func (b *Bridge) New(prefix string) interface{} {
	t, ok := b.Types[prefix]
	if !ok {
		panic("resource type not registered: " + prefix)
	}
	h := handle.New(prefix)
	r := reflect.New(t).Interface()
	handle.Set(r, h.Handle())
	return r
}

func (b *Bridge) Release(hstr string) (err error) {
	b.Lock()
	h := handle.Handle(hstr)
	b.released[h] = true
	core.Dispatch(func() {
		if err := b.Reconcile(); err != nil {
			b.Log(err)
		}
		b.Unlock()
	})
	return nil
}

func (s *Bridge) Sync(hstr string, patch map[string]interface{}, call *rpc.Call) (handle.Handle, error) {
	s.Lock()
	h := handle.Handle(hstr)

	Walk(patch, func(v, p reflect.Value, path []string) error {
		if path[len(path)-1] == "$fnptr" {
			p.SetMapIndex(reflect.ValueOf("Caller"), reflect.ValueOf(call.Caller))
		}
		return nil
	})

	v, err := s.Lookup(h)
	if err != nil {
		return h, err
	}
	if !v.IsValid() {
		v = reflect.ValueOf(s.New(h.Prefix()))
		handle.Set(v.Interface(), h.Handle())
		delete(patch, "Handle")
	}
	if err := mapstructure.Decode(patch, v.Interface()); err != nil {
		return h, err
	}
	s.Resources = append(s.Resources, v.Interface())

	core.Dispatch(func() {
		if err := s.Reconcile(); err != nil {
			s.Log(err)
		}
		s.Unlock()
	})
	return h, err
}

func (s *Bridge) Lookup(h handle.Handle) (found reflect.Value, err error) {
	for _, r := range s.Resources {
		if !handle.Has(r) {
			continue
		}
		hh := handle.Get(r)
		if hh != nil && *hh == h {
			found = reflect.ValueOf(r)
			return found, err
		}
	}
	return found, err
}

func (s *Bridge) Reconcile() error {
	for _, r := range s.Resources {
		if !handle.Has(r) {
			continue
		}
		h := handle.Get(r)
		if h != nil {
			target, exists := s.objects[*h]
			if s.released[*h] {
				// if in released but not in objects,
				// its stale state that should have been cleaned up
				// so we will ignore it here
				if !exists {
					continue
				}
				rd, ok := r.(Discarder)
				if ok && target != nil {
					if err := rd.Discard(target); err != nil {
						//delete(s.objects, *h)
						return err
					}
				}
				delete(s.objects, *h)
				continue
			}
			ra, ok := r.(Applier)
			if ok {
				var err error
				target, err = ra.Apply(target)
				if err != nil {
					return err
				}
				s.objects[*h] = target
			}
		}
	}
	return nil
}

type Applier interface {
	Apply(objc.Object) (objc.Object, error)
}

type Discarder interface {
	Discard(objc.Object) error
}