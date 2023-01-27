package wendyrpc

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Meduzz/wendy"
	"github.com/nats-io/nats.go"
)

// ServeModules serves a bunch of modules over nats.
func ServeModules(conn *nats.Conn, queue string, modules ...*wendy.Module) error {
	for _, m := range modules {
		if queue != "" {
			_, err := conn.QueueSubscribe(fmt.Sprintf("%s.*", m.Name()), queue, wrapModule(m))

			if err != nil {
				return err
			}
		} else {
			_, err := conn.Subscribe(fmt.Sprintf("%s.*", m.Name()), wrapModule(m))

			if err != nil {
				return err
			}
		}
	}

	return nil
}

// ServeMethod serves a single wendy method over nats.
func ServeMethod(conn *nats.Conn, queue, module, method string, handler wendy.Handler) error {
	if queue != "" {
		_, err := conn.QueueSubscribe(fmt.Sprintf("%s.%s", module, method), queue, wrapHandler(handler))

		if err != nil {
			return err
		}
	} else {
		_, err := conn.Subscribe(fmt.Sprintf("%s.%s", module, method), wrapHandler(handler))

		if err != nil {
			return err
		}
	}

	return nil
}

func wrapModule(module *wendy.Module) nats.MsgHandler {
	return func(msg *nats.Msg) {
		req := &wendy.Request{}
		err := json.Unmarshal(msg.Data, req)

		if err != nil {
			log.Printf("[message handler] data could not be parsed to wendy request: %s", msg.Subject)
			return
		}

		h := module.Method(req.Method)

		if h == nil {
			log.Printf("[message handler] no handler found for %s", req.Method)
			return
		}

		res := h(req)

		bs, err := json.Marshal(res)

		if err != nil {
			log.Printf("[message handler] response could not be turned in to json %s", req.Method)
			return
		}

		msg.Respond(bs)
	}
}

func wrapHandler(handler wendy.Handler) nats.MsgHandler {
	return func(msg *nats.Msg) {
		req := &wendy.Request{}
		err := json.Unmarshal(msg.Data, req)

		if err != nil {
			log.Printf("[message handler] data could not be parsed to wendy request: %s", msg.Subject)
			return
		}

		res := handler(req)

		bs, err := json.Marshal(res)

		if err != nil {
			log.Printf("[message handler] response could not be turned in to json %s", req.Method)
			return
		}

		msg.Respond(bs)
	}
}
