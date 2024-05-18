package wendyrpc

import (
	"fmt"
	"time"

	"github.com/Meduzz/rpc"
	"github.com/Meduzz/wendy"
)

type (
	WendyRpcClient struct {
		rpc     *rpc.RPC
		timeout time.Duration
	}

	Event struct {
		Body *wendy.Body
	}
)

// A building block to build wendy based clients.
func NewWendyRpcClient(rpc *rpc.RPC, timeout time.Duration) *WendyRpcClient {
	return &WendyRpcClient{rpc, timeout}
}

func (w *WendyRpcClient) Request(prefix string, req *wendy.Request) (*wendy.Response, error) {
	topic := fmt.Sprintf("%s.%s", req.Module, req.Method)

	if prefix != "" {
		topic = fmt.Sprintf("%s.%s.%s", prefix, req.Module, req.Method)
	}

	response, err := w.rpc.Request(topic, req, int(w.timeout.Seconds()))

	if err != nil {
		return nil, err
	}

	res := &wendy.Response{}
	err = response.Bind(res)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (w *WendyRpcClient) Trigger(topic string, event *Event) error {
	return w.rpc.Trigger(topic, event)
}
