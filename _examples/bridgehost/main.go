package main

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"os"

	"github.com/progrium/macbridge/res"
)

func main() {
	m, err := res.NewManager(os.Stderr)
	if err != nil {
		log.Fatal(err)
	}
	m.Start()

	data, err := ioutil.ReadFile("/Users/progrium/Source/github.com/manifold/tractor/data/icons/tractor_dark.ico")
	if err != nil {
		log.Fatal(err)
	}

	// h.Peer.Bind("Invoke", bridge.Invoke)
	// go h.Peer.Respond()

	window := res.Window{
		Title:       "Hello 1",
		Size:        res.Size{W: 480, H: 240},
		Position:    res.Point{X: 200, Y: 200},
		Closable:    true,
		Minimizable: false,
		Resizable:   false,
		Borderless:  false,
		// Image:       base64.StdEncoding.EncodeToString(data),
		// Background:   &Color{R: 0, G: 0, B: 1, A: 0.5},
	}
	if err := m.Sync(&window); err != nil {
		log.Fatal(err)
	}

	systray := res.Indicator{
		Menu: &res.Menu{
			Items: []res.MenuItem{
				{Title: "Bar", Enabled: true},
				{Title: "Foo", Enabled: true},
				{Separator: true},
				{Title: "Quit", Enabled: true},
			},
		},
		Icon: base64.StdEncoding.EncodeToString(data),
	}
	if err := m.Sync(&systray); err != nil {
		log.Fatal(err)
	}

	m.Wait()

}
