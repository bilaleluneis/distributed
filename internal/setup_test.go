/*	@author: Bilal El Uneis
	@since: August 2022
	@email: bilaleluneis@gmail.com	*/

package internal

import (
	"distributed/common"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	service := RpcNodeService{
		Nodes: make(map[common.GRPID][]RpcNode, 0),
		Ops:   make(map[common.GRPID][]FunctionalOp, 0),
	}
	if worker, err := NewWorker(service, 8080); err == nil {
		common.RegisterWorker("localhost", 8080)
		go worker.Start()
		os.Exit(m.Run())
	}
	log.Fatal("worker init failure, tests will abort")
}
