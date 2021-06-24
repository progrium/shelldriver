package shell

import (
	"github.com/progrium/shelldriver/handle"
)

type Indicator struct {
	handle.Handle

	Icon string
	Text string
	Menu *Menu
}
