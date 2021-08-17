package shell

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/mitchellh/mapstructure"
	"github.com/progrium/qtalk-go/codec"
	"github.com/progrium/qtalk-go/fn"
	"github.com/progrium/qtalk-go/mux"
	"github.com/progrium/qtalk-go/talk"
	"github.com/progrium/shelldriver/handle"
)

type Shell struct {
	Debug io.Writer
	peer  *talk.Peer
	cmd   *exec.Cmd
}

func New(sess *mux.Session) *Shell {
	var cmd *exec.Cmd
	if sess == nil {
		var err error
		bridgecmd := os.Getenv("BRIDGECMD")
		if bridgecmd == "" {
			bridgecmd = "shellbridge"
		}
		cmd = exec.Command(bridgecmd)
		wc, err := cmd.StdinPipe()
		if err != nil {
			panic(err)
		}
		rc, err := cmd.StdoutPipe()
		if err != nil {
			panic(err)
		}
		sess, err = mux.DialIO(wc, rc)
		if err != nil {
			panic(err)
		}
	}
	return &Shell{
		peer: talk.NewPeer(sess, codec.JSONCodec{}),
		cmd:  cmd,
	}
}

func (sh *Shell) Open() error {
	if sh.cmd != nil {
		sh.cmd.Stderr = sh.Debug
		if sh.Debug != nil {
			sh.cmd.Args = append(sh.cmd.Args, "-debug")
		}
		if err := sh.cmd.Start(); err != nil {
			return err
		}
	}
	go sh.peer.Respond()
	return nil
}

func (sh *Shell) Wait() error {
	if sh.cmd == nil {
		return nil
	}
	return sh.cmd.Wait()
}

func (sh *Shell) Close() error {
	if sh.cmd != nil {
		if err := sh.cmd.Process.Kill(); err != nil {
			return err
		}
	}
	return sh.peer.Close()
}

func (sh *Shell) Discard(v interface{}) error {
	if !handle.Has(v) {
		return fmt.Errorf("discard: not a resource")
	}
	h := handle.Get(v)
	if h.IsZero() {
		return fmt.Errorf("discard: cannot discard resource with zero handle")
	}
	fn.UnregisterPtrs(sh.peer.RespondMux, v)
	ctx := context.Background()
	_, err := sh.peer.Call(ctx, "Discard", []interface{}{h}, nil)
	if err == nil {
		handle.Set(v, "")
	}
	return err
}

func (sh *Shell) Sync(v interface{}) error {
	if !handle.Has(v) {
		return fmt.Errorf("sync: not a resource")
	}
	if handle.Get(v).IsZero() {
		// make sure its type is set
		handle.Set(v, "")
	}
	fn.RegisterPtrs(sh.peer.RespondMux, v)
	ctx := context.Background()
	var res interface{}
	_, err := sh.peer.Call(ctx, "Sync", []interface{}{v}, &res)
	if err != nil {
		return err
	}
	return mapstructure.Decode(res, v)
}
