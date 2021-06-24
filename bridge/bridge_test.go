package bridge

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"testing"

	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/qtalk-go/fn"
	"github.com/progrium/shelldriver/handle"
	"github.com/progrium/shelldriver/shell"
)

func init() {
	runtime.LockOSThread()
}

func TestMain(m *testing.M) {
	go func() {
		os.Exit(m.Run())
	}()
	app := cocoa.NSApp()
	app.SetActivationPolicy(cocoa.NSApplicationActivationPolicyAccessory)
	app.ActivateIgnoringOtherApps(true)
	app.Run()
}

func TestSync(t *testing.T) {
	b := New()
	defer b.Close()
	w := shell.Window{
		Title:    "Title",
		Size:     shell.Size{W: 480, H: 240},
		Position: shell.Point{X: 0, Y: 0},
		Center:   true,
	}

	if err := b.Sync(&w); err != nil {
		t.Fatal(err)
	}

	if handle.Get(w).IsZero() {
		t.Fatal("unexpected zero handle after sync")
	}

	w.Title = "New Title"
	if err := b.Sync(&w); err != nil {
		t.Fatal(err)
	}

	if w.Position.X == 0 || w.Position.Y == 0 {
		t.Fatal("unchanged position after sync with center")
	}
}

func ExampleWindow() {
	b := New()
	defer b.Close()
	w := shell.Window{
		Title:    "Title",
		Size:     shell.Size{W: 480, H: 240},
		Position: shell.Point{X: 0, Y: 0},
		Center:   true,
	}

	if err := b.Sync(&w); err != nil {
		log.Fatal(err)
	}

}

func ExampleIndicator() {
	b := New()
	i := shell.Indicator{
		Text: "Example",
		Menu: &shell.Menu{
			Items: []shell.MenuItem{
				{Title: "Test", Enabled: true, OnClick: fn.Callback(func() {
					fmt.Println("Test")
				})},
				{Separator: true},
				{Title: "Quit", Enabled: true},
			},
		},
	}

	if err := b.Sync(&i); err != nil {
		log.Fatal(err)
	}

}
