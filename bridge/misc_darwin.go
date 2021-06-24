package bridge

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
	"github.com/progrium/qtalk-go/fn"
	"github.com/progrium/shelldriver/shell"
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

func RemoteCallback(f *fn.Ptr) (objc.Object, objc.Selector) {
	ff := *f
	return core.Callback(func(o objc.Object) {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			_, err := ff.Call(ctx, nil, nil)
			if err != nil {
				fmt.Fprintf(os.Stderr, "remote callback: %v\n", err)
			}
		}()
	})
}
