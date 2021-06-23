package bridge

import (
	"github.com/progrium/macbridge/dispatch"
)

// Dispatch uses the shell API to schedule work in the main UI thread
func Dispatch(fn func() error) error {
	return dispatch.Do(dispatch.MainQueue(), fn).Wait()
}
