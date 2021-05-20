package main

import (
	"runtime"

	"github.com/progrium/macbridge/bridge"
)

const Version = "0.0.3a"

func main() {
	println("starting macbridge", Version)
	runtime.LockOSThread()
	bridge.New().Run()
}
