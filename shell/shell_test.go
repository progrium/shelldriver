package shell

import (
	"fmt"
	"log"
	"os"

	"github.com/progrium/qtalk-go/fn"
)

func ExampleWindow() {
	sh := New(nil)
	sh.Debug = os.Stderr
	if err := sh.Open(); err != nil {
		log.Fatal(err)
	}
	defer sh.Close()

	w := Window{
		Title:    "Title",
		Size:     Size{W: 480, H: 240},
		Position: Point{X: 0, Y: 0},
		Center:   true,
	}

	if err := sh.Sync(&w); err != nil {
		log.Fatal(err)
	}

	if w.Position.X == 0 {
		log.Fatal("expected position to change")
	}

	if err := sh.Discard(&w); err != nil {
		log.Fatal(err)
	}

	// // Output:
}

func ExampleIndicator() {
	sh := New(nil)
	//sh.Debug = os.Stderr
	if err := sh.Open(); err != nil {
		log.Fatal(err)
	}
	defer sh.Close()

	i := Indicator{
		Text: "Example",
		Menu: &Menu{
			Items: []MenuItem{
				{Title: "Test", Enabled: true, OnClick: fn.Callback(func() {
					fmt.Fprintln(os.Stderr, "Test")
				})},
				{Separator: true},
				{Title: "Quit", Enabled: true},
			},
		},
	}

	if err := sh.Sync(&i); err != nil {
		log.Fatal(err)
	}

	// // Output:
}
