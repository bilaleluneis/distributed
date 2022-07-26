package internal

import (
	"distributed/common"
	"fmt"
	"log"
	"net"
	"net/rpc"
)

func InitServer(port int) error {
	if err := rpc.Register(&RpcNode{}); err != nil {
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

func GetTcpClient(distn common.Location) (*rpc.Client, error) {
	address := fmt.Sprintf("%s:%d", distn.HostName, distn.Port)
	client, err := rpc.Dial("tcp", address)
	if err != nil {
		log.Print("client connection failed")
		return nil, err
	}
	return client, nil
}
