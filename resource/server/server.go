package server

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
	"github.com/progrium/macbridge/bridge"
	"github.com/progrium/macbridge/handle"
	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
)

type Server struct {
	Resources []interface{}
	Types     map[string]reflect.Type

	objects  map[handle.Handle]objc.Object
	released map[handle.Handle]bool

	sync.Mutex
}

func New() *Server {
	b := &Server{
		Types:    make(map[string]reflect.Type),
		objects:  make(map[handle.Handle]objc.Object),
		released: make(map[handle.Handle]bool),
	}
	b.Register(&bridge.Indicator{})
	b.Register(&bridge.Menu{})
	b.Register(&bridge.Window{})
	return b
}

func (b *Server) Run() {
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
	app.SetActivationPolicy(cocoa.NSApplicationActivationPolicyAccessory)
	app.ActivateIgnoringOtherApps(true)
	app.Run()
}

func (b *Server) Log(v ...interface{}) {
	fmt.Fprintln(os.Stderr, v...)
}

func (b *Server) Register(v interface{}) {
	h := handle.Get(v)
	b.Types[h.Type()] = reflect.TypeOf(v).Elem()
}

func (b *Server) New(typ string) interface{} {
	t, ok := b.Types[typ]
	if !ok {
		panic("resource type not registered: " + typ)
	}
	r := reflect.New(t).Interface()
	h := handle.NewFor(r)
	handle.Set(r, h.Handle())
	return r
}

func (b *Server) Release(hstr string) (err error) {
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

func (s *Server) Sync(res map[string]interface{}, call *rpc.Call) (handle.Handle, error) {
	s.Lock()
	h := handle.Handle(res["Handle"].(string))
	if h.Type() == "" {
		panic("invalid handle")
	}

	Walk(res, func(v, p reflect.Value, path []string) error {
		if path[len(path)-1] == "$fnptr" {
			p.SetMapIndex(reflect.ValueOf("Caller"), reflect.ValueOf(call.Caller))
		}
		return nil
	})

	v, err := s.Lookup(h)
	if err != nil {
		return h, err
	}
	var r interface{}
	if !v.IsValid() {
		v = reflect.ValueOf(s.New(h.Type()))
		rv, ok := v.Interface().(handle.Resourcer)
		if !ok {
			panic("sync: not a resource")
		}
		_, r = rv.Resource()
		handle.Set(r, handle.NewFor(r).Handle())
		delete(res, "Handle")
	}
	if err := mapstructure.Decode(res, r); err != nil {
		return h, err
	}
	s.Resources = append(s.Resources, v.Interface())

	println(fmt.Sprintf("%#v", s.Resources[0]))

	core.Dispatch(func() {
		if err := s.Reconcile(); err != nil {
			s.Log(err)
		}
		s.Unlock()
	})
	return h, err
}

func (s *Server) Lookup(h handle.Handle) (found reflect.Value, err error) {
	for _, r := range s.Resources {
		if !handle.Has(r) {
			continue
		}
		hh := handle.Get(r)
		if !hh.Unset() && hh == h {
			found = reflect.ValueOf(r)
			return found, err
		}
	}
	return found, err
}

func (s *Server) Reconcile() error {
	for _, r := range s.Resources {
		if !handle.Has(r) {
			continue
		}
		h := handle.Get(r)
		if !h.Unset() {
			target, exists := s.objects[h]
			if s.released[h] {
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
				delete(s.objects, h)
				continue
			}
			ra, ok := r.(Applier)
			if ok {
				var err error
				target, err = ra.Apply(target)
				if err != nil {
					return err
				}
				s.objects[h] = target
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
