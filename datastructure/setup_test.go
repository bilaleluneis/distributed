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
		if worker, err := internal.NewWorker(port); err == nil {
			common.RegisterWorker("localhost", port)
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
