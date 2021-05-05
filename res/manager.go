package res

import (
	"context"
	"io"
	"os"
	"os/exec"

	"github.com/manifold/qtalk/golang/mux"
	"github.com/manifold/qtalk/golang/rpc"
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
		bridgecmd = "macbridge"
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
