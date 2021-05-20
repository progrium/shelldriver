package shell

import (
	"github.com/progrium/macbridge/handle"
)

type Indicator struct {
	handle.Handle

	Icon string
	Text string
	Menu *Menu
}
