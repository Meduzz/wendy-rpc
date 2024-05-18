package main

import (
	"encoding/json"
	"net/http"

	"github.com/Meduzz/rpc"
	"github.com/Meduzz/wendy"
	wendyrpc "github.com/Meduzz/wendy-rpc"
	"github.com/nats-io/nats.go"
)

func main() {
	conn, err := nats.Connect("nats://localhost:4222")

	if err != nil {
		panic(err)
	}

	srv := rpc.NewRpc(conn)
	proxy := wendyrpc.WendyProxy(srv, "")

	http.HandleFunc("/api", handle(proxy))
	err = http.ListenAndServe(":8080", http.DefaultServeMux)

	if err != nil {
		panic(err)
	}
}

func handle(proxy wendy.Wendy) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		wendyReq := &wendy.Request{}
		decoder := json.NewDecoder(req.Body)

		err := decoder.Decode(wendyReq)

		if err != nil {
			res.WriteHeader(400)
			return
		}

		wendyRes := proxy.Handle(req.Context(), wendyReq)

		if wendyRes.Headers != nil {
			for k, v := range wendyRes.Headers {
				res.Header().Add(k, v)
			}
		}

		if wendyRes.Body != nil {
			res.Header().Add("Content-Type", wendyRes.Body.Type)
			res.WriteHeader(wendyRes.Code)
			res.Write(wendyRes.Body.Data)
		} else {
			res.WriteHeader(wendyRes.Code)
		}
	}
}
