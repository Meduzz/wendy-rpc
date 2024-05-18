package wendyrpc

import (
	"context"
	"fmt"
	"log"

	"github.com/Meduzz/rpc"
	"github.com/Meduzz/wendy"
	"github.com/nats-io/nats.go"
)

type (
	wendyProxy struct {
		prefix string
		srv    *rpc.RPC
	}
)

func (w *wendyProxy) Handle(ctx context.Context, req *wendy.Request) *wendy.Response {
	topic := fmt.Sprintf("%s.%s", req.Module, req.Method)

	if w.prefix != "" {
		topic = fmt.Sprintf("%s.%s.%s", w.prefix, req.Module, req.Method)
	}

	resCtx, err := w.srv.RequestContext(ctx, topic, req)

	if err != nil {
		if err == nats.ErrTimeout {
			log.Printf("Request to %s timed out\n", topic)
			return &wendy.Response{
				Code: 503,
			}
		} else {
			log.Printf("Request to %s threw error: %v\n", topic, err)
			return wendy.Error(nil)
		}
	}

	res := &wendy.Response{}
	err = resCtx.Bind(res)

	if err != nil {
		log.Printf("Parsing response threw error: %v\n", err)
		return wendy.Error(nil)
	}

	return res
}

// WendyProxy implements Wendy interface but proxies requests over nats
func WendyProxy(srv *rpc.RPC, prefix string) wendy.Wendy {
	return &wendyProxy{prefix, srv}
}
