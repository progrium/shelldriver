package res

import (
	"github.com/progrium/macbridge/handle"
)

type Menu struct {
	*handle.Handle `prefix:"men"`

	Icon    string
	Title   string
	Tooltip string
	Items   []MenuItem
}

type MenuItem struct {
	*handle.Handle `prefix:"mit"`

	Title     string
	Icon      string
	Tooltip   string
	Separator bool
	Enabled   bool
	Checked   bool

	//OnClick *rpc_FuncExport
	// TODO: submenus
}
