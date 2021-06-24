package bridge

import (
	"encoding/base64"

	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/webkit"
	"github.com/progrium/shelldriver/shell"
)

func init() {
	register(&Window{})
}

type Window struct {
	shell.Window `mapstructure:",squash"`

	target  *cocoa.NSWindow
	webview *webkit.WKWebView
	image   *cocoa.NSImage
}

func (w *Window) Resource() interface{} {
	return &w.Window
}

func (w *Window) Discard() error {
	// TODO: image, webview
	w.target.Close()
	w.target.Release()
	return nil
}

func (w *Window) Apply() error {
	frame := core.Rect(w.Position.X, w.Position.Y, w.Size.W, w.Size.H)

	if w.target == nil {
		obj := cocoa.NSWindow_Init(frame, cocoa.NSTitledWindowMask, cocoa.NSBackingStoreBuffered, false)
		w.target = &obj
		w.target.Retain()
		w.target.MakeKeyAndOrderFront(nil)
		if w.Center {
			screenRect := cocoa.NSScreen_Main().Frame()
			w.Position.X = (screenRect.Size.Width / 2) - (w.Size.W / 2)
			w.Position.Y = (screenRect.Size.Height / 2) - (w.Size.H / 2)
			frame = core.Rect(w.Position.X, w.Position.Y, w.Size.W, w.Size.H)
		}
	}

	// URL
	if w.URL != "" && w.webview == nil {
		config := webkit.WKWebViewConfiguration_New()
		config.Preferences().SetValueForKey(core.True, core.String("developerExtrasEnabled"))

		wv := webkit.WKWebView_Init(core.Rect(0, 0, w.Size.W, w.Size.H), config)
		w.webview = &wv

		req := core.NSURLRequest_Init(core.URL(w.URL))
		wv.LoadRequest(req)
	}

	// Image
	if w.Image != "" && w.image == nil {
		b, err := base64.StdEncoding.DecodeString(w.Image)
		if err != nil {
			return err
		}
		data := core.NSData_WithBytes(b, uint64(len(b)))
		image := cocoa.NSImage_InitWithData(data)
		w.image = &image
	}

	// Closabe, Minimizable, Borderless, Resizable
	mask := cocoa.NSTitledWindowMask
	needsTitleBar := w.Closable || w.Minimizable
	if w.Borderless {
		if !needsTitleBar {
			mask = cocoa.NSBorderlessWindowMask
		}
		mask = mask | cocoa.NSFullSizeContentViewWindowMask
	}
	if w.Closable {
		mask = mask | cocoa.NSClosableWindowMask
	}
	if w.Minimizable {
		mask = mask | cocoa.NSMiniaturizableWindowMask
	}
	if w.Resizable {
		mask = mask | cocoa.NSResizableWindowMask
	}
	w.target.SetStyleMask(mask)

	// Title
	if w.Title != "" {
		w.target.SetTitle(w.Title)
	} else {
		w.target.SetMovableByWindowBackground(true)
		w.target.SetTitlebarAppearsTransparent(true)
	}

	// Borderless, CornerRadius, Background
	if w.Borderless && w.CornerRadius > 0 {
		w.target.SetBackgroundColor(cocoa.NSColor_Clear())
		w.target.SetOpaque(false)
		v := cocoa.NSView_Init(core.Rect(0, 0, 0, 0))
		if w.Background != nil {
			v.SetBackgroundColor(NSColor(w.Background))
		}
		v.SetWantsLayer(true)
		v.Layer().SetCornerRadius(w.CornerRadius)

		if w.webview != nil {
			v.AddSubviewPositionedRelativeTo(*w.webview, cocoa.NSWindowAbove, nil)
		}
		w.target.SetContentView(v)
	} else {
		if w.Background != nil {
			w.target.SetBackgroundColor(NSColor(w.Background))
			w.target.SetOpaque(w.Background.A == 1)
		}
		if w.webview != nil {
			w.target.SetContentView(*w.webview)
		}
	}

	// Background
	if w.webview != nil && w.Background != nil && w.Background.A == 0 {
		w.webview.SetOpaque(false)
		w.webview.SetBackgroundColor(cocoa.NSColor_Clear())
		w.webview.SetValueForKey(core.False, core.String("drawsBackground"))
	}

	// Image
	if w.image != nil {
		w.target.ContentView().SetWantsLayer(true)
		w.target.ContentView().Layer().SetContents(w.image)
	}

	// AlwaysOnTop
	if w.AlwaysOnTop {
		w.target.SetLevel(cocoa.NSMainMenuWindowLevel)
	}

	// IgnoreMouse
	if w.IgnoreMouse {
		w.target.SetIgnoresMouseEvents(true)
	}

	w.target.SetFrameDisplay(frame, true)

	return nil
}
