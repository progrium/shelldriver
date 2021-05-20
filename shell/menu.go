package shell

import (
	"github.com/progrium/macbridge/handle"
)

type Menu struct {
	handle.Handle

	Icon    string
	Title   string
	Tooltip string
	Items   []MenuItem
}

type MenuItem struct {
	handle.Handle

	Title     string
	Icon      string
	Tooltip   string
	Separator bool
	Enabled   bool
	Checked   bool

	//OnClick *rpc_FuncExport
	// TODO: submenus
}
