# macbridge
Common language bridge for common macOS resources.

A work in progress recently split out of [progrium/macdriver](https://github.com/progrium/macdriver). Original section of readme:

## Bridge System
Lastly, a common case for this toolkit is not just building full native apps, but integrating Go applications
with Mac systems, like windows, native menus, status icons (systray), etc.
One-off libraries for some of these exist, but besides often limiting what you can do, 
they're also just not composable. They all want to own the main thread!

For this and other reasons, we often run the above kind of code in a separate process altogether from our
Go application. This might seem like a step backwards, but it is safer and more robust in a way. 

The `bridge` package takes advantage of this situation to create a higher-level abstraction more aligned with a potential 
cross-platform toolkit. You can declaratively describe and modify structs that can be copied to the bridge process and applied to the Objective-C
objects in a manner similar to configuration management:

```go
package main 

import (
	"os"

	"github.com/progrium/macdriver/bridge"
)

func main() {
	// start a bridge subprocess
	host := bridge.NewHost(os.Stderr)
	go host.Run()

	// create a window
	window := bridge.Window{
		Title:       "My Title",
		Size:        bridge.Size{W: 480, H: 240},
		Position:    bridge.Point{X: 200, Y: 200},
		Closable:    true,
		Minimizable: false,
		Resizable:   false,
		Borderless:  false,
		AlwaysOnTop: true,
		Background:   &bridge.Color{R: 1, G: 1, B: 1, A: 0.5},
	}
	host.Sync(&window)

	// change its title
	window.Title = "My New Title"
	host.Sync(&window)

	// destroy the window
	host.Release(&window)
}

```
This is the most WIP part of the project, but once developed further we can take this API and build a bridge
system with the same resources for Windows and Linux, making a cross-platform OS "driver". We'll see.

* Current bridge types available:
  * Window
  * StatusItem (systray)
  * Menu