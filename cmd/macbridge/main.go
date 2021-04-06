package main

import (
	"runtime"

	"github.com/progrium/macbridge/pkg/bridge"
)

const Version = "0.0.2"

func main() {
	println("starting macbridge", Version)
	runtime.LockOSThread()
	bridge.Run()
}
