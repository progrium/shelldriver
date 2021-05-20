package main

import (
	"runtime"

	"github.com/progrium/macbridge/resource/server"
)

const Version = "0.0.4"

func main() {
	println("starting sdbridge", Version)
	runtime.LockOSThread()
	server.New().Run()
}
