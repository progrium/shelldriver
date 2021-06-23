package bridge

import (
	"errors"

	"github.com/progrium/macbridge/handle"
	"github.com/progrium/qtalk-go/codec"
	"github.com/progrium/qtalk-go/fn"
	"github.com/progrium/qtalk-go/rpc"
)

func NewServer() *rpc.Server {
	return &rpc.Server{
		Codec: codec.JSONCodec{},
		Handler: fn.HandlerFrom(&server{
			bridge: New(),
		}),
	}
}

type server struct {
	bridge *Bridge
}

func (s *server) Discard(h string) error {
	return s.bridge.Discard(handle.Handle(h))
}

func (s *server) Sync(res map[string]interface{}, call *rpc.Call) (interface{}, error) {
	hstr, ok := res["Handle"].(string)
	if !ok {
		return nil, errors.New("no Handle string")
	}
	fn.SetCallers(&res, call.Caller)
	return s.bridge.Apply(handle.Handle(hstr), res)
}
