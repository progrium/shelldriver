package res

import (
	"github.com/progrium/macbridge/handle"
)

type Indicator struct {
	*handle.Handle `prefix:"ind"`

	Icon string
	Text string
	Menu *Menu
}
