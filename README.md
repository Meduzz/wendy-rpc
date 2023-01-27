# wendy-rpc
rpc bindings for wendy

## Use

### Server

To build a wendy server on the RPC transport, do something like this:

    // where conn is your nats connection, and service.ServiceModule() returns a wendy module
	err = wendyrpc.ServeModules(conn, "workgroup1", service.ServiceModule())

For more inspiration of this, see the example/server.

### Proxy
To get a RPC proxy that will forward wendy requests over nats, do like this:

    // conn is your nats connection
	srv := rpc.NewRpc(conn)
	proxy := wendyrpc.WendyProxy(srv)

    // proxy implements wendy.Wendy and can be called like this:
    wendyRes := proxy.Handle(req.Context(), wendyReq)

For inspiration of how to use the proxy, see the example/proxy.