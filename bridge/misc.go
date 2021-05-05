package bridge

import (
	"reflect"

	"github.com/progrium/macbridge/res"
	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/core"
)

func prefix(v interface{}) string {
	rv := reflect.Indirect(reflect.ValueOf(v))
	return rv.Type().Field(0).Tag.Get("prefix")
}

func NSPoint(p *res.Point) core.NSPoint {
	return core.NSPoint{X: p.X, Y: p.Y}
}

func NSSize(s *res.Size) core.NSSize {
	return core.NSSize{Width: s.W, Height: s.H}
}

func NSColor(c *res.Color) cocoa.NSColor {
	return cocoa.NSColor_Init(c.R, c.G, c.B, c.A)
}
