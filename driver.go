package distributed

import (
	"distributed/internal"
	"fmt"
	"log"
	"net"
	"net/rpc"
)

func AsServer(port int) error {
	if err := rpc.Register(&internal.RpcNode{}); err != nil {
		return err
	}
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Print("server failed to initialize")
		return err
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print("server failed to accept connection")
			return err
		}
		go rpc.ServeConn(conn)
	}
}

var client *rpc.Client

func AsClient(host string, port int) error {
	if client == nil {
		c, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
		if err != nil {
			log.Print("client connection failed")
			return err
		}
		client = c
	}
	return nil
}
