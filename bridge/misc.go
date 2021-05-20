package bridge

import (
	"github.com/progrium/macbridge/shell"
	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/core"
)

func NSPoint(p *shell.Point) core.NSPoint {
	return core.NSPoint{X: p.X, Y: p.Y}
}

func NSSize(s *shell.Size) core.NSSize {
	return core.NSSize{Width: s.W, Height: s.H}
}

func NSColor(c *shell.Color) cocoa.NSColor {
	return cocoa.NSColor_Init(c.R, c.G, c.B, c.A)
}
