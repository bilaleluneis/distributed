package internal

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
)

func InitServer(port int) error {
	if err := rpc.Register(&RpcNodeService{}); err != nil {
		return err
	}
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
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
