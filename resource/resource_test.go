package resource

import (
	"context"
	"io"
	"log"
	"testing"

	"github.com/manifold/qtalk/golang/mux"
	"github.com/manifold/qtalk/golang/rpc"
	"github.com/progrium/macbridge/handle"
)

func mockServer() (io.Reader, io.WriteCloser) {
	r, sw := io.Pipe()
	sr, w := io.Pipe()
	go func() {
		session := mux.NewSession(context.Background(), struct {
			io.ReadCloser
			io.Writer
		}{sr, sw})
		peer := rpc.NewPeer(session, rpc.JSONCodec{})
		peer.Bind("Sync", rpc.HandlerFunc(func(r rpc.Responder, c *rpc.Call) {
			var res map[string]interface{}
			if err := c.Decode(&res); err != nil {
				log.Fatal(err)
			}
			r.Return(handle.New("test", "123"))
		}))
		peer.Bind("Release", rpc.HandlerFunc(func(r rpc.Responder, c *rpc.Call) {
			r.Return(nil)
		}))
		peer.Respond()
	}()
	return r, w
}

func TestSyncRelease(t *testing.T) {
	r, w := mockServer()

	pipe := struct {
		io.WriteCloser
		io.Reader
	}{w, r}
	session := mux.NewSession(context.Background(), pipe)
	manager := &Manager{
		Pipe: pipe,
		Peer: rpc.NewPeer(session, rpc.JSONCodec{}),
	}

	m := &Menu{
		Icon:  "icon",
		Title: "title",
	}
	if !m.Handle.Unset() {
		t.Fatal("menu handle should be unset")
	}

	if err := manager.Sync(m); err != nil {
		t.Fatal(err)
	}

	if m.Unset() {
		t.Fatal("menu handle should not be unset")
	}

	if err := manager.Release(m); err != nil {
		t.Fatal(err)
	}
	if !m.Handle.Unset() {
		t.Fatal("menu handle should be unset")
	}

}
