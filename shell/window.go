package shell

import "github.com/progrium/macbridge/handle"

type Window struct {
	handle.Handle

	Title        string
	Position     Point
	Size         Size
	Closable     bool
	Minimizable  bool
	Resizable    bool
	Background   *Color
	Borderless   bool
	CornerRadius float64
	AlwaysOnTop  bool
	IgnoreMouse  bool
	Center       bool
	URL          string
	Image        string
}
