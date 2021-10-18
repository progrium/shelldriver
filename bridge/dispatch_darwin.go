package bridge

import (
	"github.com/progrium/macdriver/core"
)

// Dispatch uses the shell API to schedule work in the main UI thread
func Dispatch(fn func() error) error {
	errCh := make(chan error, 1)
	core.Dispatch(func() {
		errCh <- fn()
	})
	return <-errCh
}
