package datastructure

import (
	"distributed/common"
	"distributed/internal"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	ports := []int{8080, 8081, 8083}
	for _, port := range ports {
		service := internal.RpcNodeService{
			Nodes: make(map[common.GRPID][]internal.RpcNode, 0),
			Ops:   make(map[common.GRPID][]internal.FunctionalOp, 0),
		}
		if worker, err := internal.NewWorker(service, port); err == nil {
			common.RegisterWorker("localhost", port)
			log.Printf("Test Worker Started on Port %d", port)
			go worker.Start()
		}
	}

	// all workers failed initialization.. can't proceed with tests
	if len(common.GetAvailRegWorkers()) == 0 {
		log.Fatal("non of workers started up")
	}

	// we got here, then we have workers and we can run tests
	os.Exit(m.Run())
}
