package main

import (
	"runtime"

	"github.com/progrium/macbridge/pkg/bridge"
)

func main() {
	runtime.LockOSThread()
	bridge.Run()
}
