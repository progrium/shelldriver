package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/progrium/qtalk-go/transport/qmux"
	"github.com/progrium/shelldriver/bridge"
)

const Version = "0.1.0"

func init() {
	runtime.LockOSThread()
}

func main() {
	flagDebug := flag.Bool("debug", false, "debug mode")
	flag.Parse()

	if *flagDebug {
		fmt.Fprintf(os.Stderr, "shellbridge %s\n", Version)
	}

	sess, err := qmux.DialIO(os.Stdout, os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	srv := bridge.NewServer()
	go srv.Respond(sess)

	bridge.Main()
}
