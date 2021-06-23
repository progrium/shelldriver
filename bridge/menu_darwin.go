package bridge

import (
	"encoding/base64"

	"github.com/progrium/macbridge/shell"
	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
)

func init() {
	register(&Menu{})
}

type Menu struct {
	shell.Menu `mapstructure:",squash"`

	target *cocoa.NSMenu
}

func (m *Menu) Resource() interface{} {
	return &m.Menu
}

func (m *Menu) Discard() error {
	m.target.Release()
	return nil
}

func (m *Menu) Apply() error {
	obj := cocoa.NSMenu_New()
	if m.target != nil {
		m.target.Release()
	}
	m.target = &obj
	m.target.SetAutoenablesItems(true)
	for _, i := range m.Items {
		m.target.AddItem(NSMenuItem(i))
	}

	// TODO: Tooltip
	// TODO: Icon
	// TODO: Title

	return nil
}

func NSMenuItem(i shell.MenuItem) cocoa.NSMenuItem {
	// Separator
	if i.Separator {
		return cocoa.NSMenuItem_Separator()
	}

	obj := cocoa.NSMenuItem_New()

	// Title
	obj.SetTitle(i.Title)

	// Enabled
	obj.SetEnabled(i.Enabled)

	// Tooltip
	obj.SetToolTip(i.Tooltip)

	// Checked
	if i.Checked {
		obj.SetState(cocoa.NSControlStateValueOn)
	}

	// Icon
	if i.Icon != "" {
		b, err := base64.StdEncoding.DecodeString(i.Icon)
		if err == nil {
			data := core.NSData_WithBytes(b, uint64(len(b)))
			img := cocoa.NSImage_InitWithData(data)
			img.SetSize(core.Size(16, 16))
			obj.SetImage(img)
		}
	}

	// Quit default action
	if i.Title == "Quit" {
		obj.SetTarget(cocoa.NSApp())
		obj.SetAction(objc.Sel("terminate:"))
	}

	// OnClick action
	if i.OnClick != nil {
		t, sel := RemoteCallback(i.OnClick)
		obj.SetTarget(t)
		obj.SetAction(sel)
	}

	// SubItems
	if len(i.SubItems) > 0 {
		sub := cocoa.NSMenu_New()
		sub.SetAutoenablesItems(true)
		for _, i := range i.SubItems {
			sub.AddItem(NSMenuItem(i))
		}
		obj.SetSubmenu(sub)
	}

	return obj
}
