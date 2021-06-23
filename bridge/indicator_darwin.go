package bridge

import (
	"encoding/base64"

	"github.com/progrium/macbridge/shell"
	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/core"
)

func init() {
	register(&Indicator{})
}

type Indicator struct {
	shell.Indicator `mapstructure:",squash"`

	target *cocoa.NSStatusItem
	menu   *Menu
}

func (i *Indicator) Resource() interface{} {
	return &i.Indicator
}

func (i *Indicator) Discard() error {
	if i.menu != nil {
		if err := i.menu.Discard(); err != nil {
			return err
		}
	}
	i.target.Release()
	return nil
}

func (i *Indicator) Apply() error {
	if i.target == nil {
		obj := cocoa.NSStatusBar_System().StatusItemWithLength(cocoa.NSVariableStatusItemLength)
		i.target = &obj
		i.target.Retain()
	}

	// Text
	i.target.Button().SetTitle(i.Text)

	// Icon
	if i.Icon != "" {
		b, err := base64.StdEncoding.DecodeString(i.Icon)
		if err != nil {
			return err
		}
		data := core.NSData_WithBytes(b, uint64(len(b)))
		image := cocoa.NSImage_InitWithData(data)
		image.SetSize(core.Size(16.0, 16.0))
		image.SetTemplate(true)
		i.target.Button().SetImage(image)
		if i.Text != "" {
			i.target.Button().SetImagePosition(cocoa.NSImageLeft)
		} else {
			i.target.Button().SetImagePosition(cocoa.NSImageOnly)
		}

	}

	// Menu
	if i.Menu != nil {
		if i.menu == nil {
			i.menu = &Menu{}
		}
		i.menu.Menu = *i.Menu
		if err := i.menu.Apply(); err != nil {
			return err
		}
		i.target.SetMenu(*i.menu.target)
	}

	return nil
}
