package shell

import (
	"github.com/progrium/macbridge/handle"
	"github.com/progrium/qtalk-go/fn"
)

type Menu struct {
	handle.Handle

	Icon    string
	Title   string
	Tooltip string
	Items   []MenuItem
}

type MenuItem struct {
	Title     string
	Icon      string
	Tooltip   string
	Separator bool
	Enabled   bool
	Checked   bool

	OnClick  *fn.Ptr
	SubItems []MenuItem
}
