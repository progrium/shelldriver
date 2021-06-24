package main

import (
	"log"
	"os"

	"github.com/progrium/qtalk-go/fn"
	"github.com/progrium/shelldriver/shell"
)

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	sh := shell.New(nil)
	sh.Debug = os.Stderr
	must(sh.Open())
	defer sh.Close()

	w := shell.Window{
		Title:    "Demo",
		Size:     shell.Size{W: 480, H: 240},
		Position: shell.Point{X: 0, Y: 0},
		Center:   true,
	}
	must(sh.Sync(&w))

	i := shell.Indicator{
		Text: "🚜",
		Menu: &shell.Menu{
			Items: []shell.MenuItem{
				{Title: "Window", Enabled: true, SubItems: []shell.MenuItem{
					{Title: "Step", Enabled: true, OnClick: fn.Callback(func() {
						w.Position.X += 20
						w.Position.Y += 20
						must(sh.Sync(&w))
					})},
					// TODO: Always On Top check/state
				}},
				{Separator: true},
				{Title: "Quit", Enabled: true},
			},
		},
	}
	must(sh.Sync(&i))

	sh.Wait()
}
