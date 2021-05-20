package resource

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/manifold/qtalk/golang/mux"
	"github.com/manifold/qtalk/golang/rpc"
	"github.com/progrium/macbridge/handle"
)

type Manager struct {
	Cmd  *exec.Cmd
	Pipe io.ReadWriteCloser
	Peer *rpc.Peer
}

func (m *Manager) Start() error {
	return m.Cmd.Start()
}

func (m *Manager) Wait() error {
	return m.Cmd.Wait()
}

func NewManager(stderr io.Writer) (*Manager, error) {
	bridgecmd := os.Getenv("BRIDGECMD")
	if bridgecmd == "" {
		bridgecmd = "sdbridge"
	}
	cmd := exec.Command(bridgecmd)
	cmd.Stderr = stderr
	wc, inErr := cmd.StdinPipe()
	if inErr != nil {
		return nil, inErr
	}
	rc, outErr := cmd.StdoutPipe()
	if outErr != nil {
		return nil, outErr
	}
	pipe := struct {
		io.WriteCloser
		io.Reader
	}{wc, rc}
	session := mux.NewSession(context.Background(), pipe)
	return &Manager{Cmd: cmd, Pipe: pipe, Peer: rpc.NewPeer(session, rpc.JSONCodec{})}, nil
}

func (m *Manager) Sync(v interface{}) error {
	if !handle.Has(v) {
		return fmt.Errorf("not a resource")
	}
	handle.Set(v, handle.Get(v).Handle())
	var h string
	_, err := m.Peer.Call("Sync", []interface{}{v}, &h)
	handle.Set(v, h)
	return err
}

func (m *Manager) Release(v interface{}) error {
	if !handle.Has(v) {
		return fmt.Errorf("not a resource")
	}
	h := handle.Get(v)
	if h.Unset() {
		return fmt.Errorf("unable to release an unset resource handle")
	}
	_, err := m.Peer.Call("Release", h, nil)
	if err == nil {
		handle.Set(v, "")
	}
	return err
}
