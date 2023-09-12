package main

import (
	"os"
	"os/signal"

	wendyrpc "github.com/Meduzz/wendy-rpc"
	"github.com/Meduzz/wendy/example/service"
	"github.com/nats-io/nats.go"
)

func main() {
	conn, err := nats.Connect("nats://localhost:4222")

	if err != nil {
		panic(err)
	}

	err = wendyrpc.ServeModules(conn, "workgroup1", "example", service.ServiceModule())

	if err != nil {
		panic(err)
	}

	defer conn.Close()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
