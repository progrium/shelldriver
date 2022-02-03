package bridge

import (
	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
)

func Main() {
	app := cocoa.NSApp()
	nsbundleMain := cocoa.NSBundle_Main()
	nsbundle := nsbundleMain.Class()
	nsbundle.AddMethod("__bundleIdentifier", func(self objc.Object) objc.Object {
		if self.Pointer() == nsbundleMain.Pointer() {
			return core.String("com.progrium.shellbridge")
		}
		// After the swizzle this will point to the original method, and return the
		// original bundle identifier.
		return self.Send("__bundleIdentifier")
	})
	nsbundle.Swizzle("bundleIdentifier", "__bundleIdentifier")
	app.SetActivationPolicy(cocoa.NSApplicationActivationPolicyAccessory)
	app.ActivateIgnoringOtherApps(true)
	app.Run()
}
