/*	@author: Bilal El Uneis
	@since: August 2022
	@email: bilaleluneis@gmail.com	*/

package internal

import (
	"distributed/common"
	"fmt"
	"log"
	"net"
	"net/rpc"
)

type Worker struct {
	handler  *rpc.Server
	listener net.Listener
}

func NewWorker(atPort int) (Worker, error) {
	var err error
	var worker Worker
	service := RpcNodeService{make(map[common.GRPID][]RpcNode, 0)}
	worker.handler = rpc.NewServer()
	if err = worker.handler.Register(&service); err != nil {
		return worker, common.RpcServiceRegErr
	}
	address := fmt.Sprintf("0.0.0.0:%d", atPort)
	if worker.listener, err = net.Listen("tcp", address); err != nil {
		return worker, common.InitWorkerFailed
	}
	return worker, nil
}

func (w Worker) Start() {
	for {
		conn, err := w.listener.Accept()
		if err != nil {
			log.Printf("failed to process work request due: %s", err)
			continue
		}
		go w.handler.ServeConn(conn)
	}
}
