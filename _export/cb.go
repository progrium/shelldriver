package export

import (
	"fmt"
	"os"

	"github.com/manifold/qtalk/golang/rpc"
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
	"github.com/rs/xid"
)

type rpc_FuncExport struct {
	Ptr    string `json:"$fnptr" mapstructure:"$fnptr"`
	Caller rpc.Caller
	fn     interface{}
}

func (e *rpc_FuncExport) Call(args, reply interface{}) error {
	_, err := e.Caller.Call("Invoke", e.Ptr, reply)
	return err
}

func (e *rpc_FuncExport) Callback() (objc.Object, objc.Selector) {
	ee := *e
	return core.Callback(func(o objc.Object) {
		err := ee.Call(nil, nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "callback: %v\n", err)
		}
	})
}

var exportedFuncs map[string]rpc_FuncExport

func ExportFunc(fn interface{}) *rpc_FuncExport {
	if exportedFuncs == nil {
		exportedFuncs = make(map[string]rpc_FuncExport)
	}
	id := xid.New().String()
	ef := rpc_FuncExport{
		Ptr: id,
		fn:  fn,
	}
	exportedFuncs[id] = ef
	return &ef
}
