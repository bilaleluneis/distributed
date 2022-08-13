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

// moved the registration of RPC service to init to ensure
// one time registration upon inclusion of module regardless
// of how many worker instances created.. this makes it possible
// to create multiple workers in unit tests.
func init() {
	if err := rpc.Register(&RpcNodeService{}); err != nil {
		panic(common.RpcServiceRegErr)
	}
}

func InitWorker(port int) (net.Listener, error) {
	listnr, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		log.Printf("worker failed to start on port %d", port)
		return nil, err
	}
	return listnr, nil
}

func ProcessWorkRequest(l net.Listener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("failed to process work request due: %s", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
