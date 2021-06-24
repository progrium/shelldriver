package bridge

import (
	"context"
	"io"
	"testing"

	"github.com/progrium/qtalk-go/codec"
	"github.com/progrium/qtalk-go/rpc"
	"github.com/progrium/qtalk-go/transport/qmux"
	"github.com/progrium/shelldriver/handle"
	"github.com/progrium/shelldriver/shell"
)

func TestServer(t *testing.T) {
	ar, bw := io.Pipe()
	br, aw := io.Pipe()
	sessA, _ := qmux.DialIO(aw, ar)
	sessB, _ := qmux.DialIO(bw, br)

	srv := NewServer()
	go srv.Respond(sessA)

	client := rpc.NewClient(sessB, codec.JSONCodec{})

	w := shell.Window{
		Title:    "Title",
		Size:     shell.Size{W: 480, H: 240},
		Position: shell.Point{X: 0, Y: 0},
		Center:   true,
	}
	handle.Set(&w, "")

	var res map[string]interface{}
	_, err := client.Call(context.Background(), "Sync", []interface{}{w}, &res)
	if err != nil {
		t.Fatal(err)
	}

	if res["Handle"].(string) == "" {
		t.Fatal("expected Handle to be set")
	}
}
