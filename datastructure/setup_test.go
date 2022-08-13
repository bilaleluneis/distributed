package datastructure

import (
	"distributed/common"
	"distributed/internal"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	// not all those might succeed initialization
	workers := []common.Worker{
		{Host: "localhost", Port: 8080},
		{Host: "localhost", Port: 8081},
		{Host: "localhost", Port: 8083},
	}

	startedWorkers := make([]common.Worker, 0)
	for _, worker := range workers {
		if l, err := internal.InitWorker(worker.Port); err == nil {
			startedWorkers = append(startedWorkers, worker)
			go internal.ProcessWorkRequest(l)
			log.Printf("worker at port %d started successfully", worker.Port)
		}
	}

	// all workers failed initialization.. can't proceed with tests
	if len(startedWorkers) == 0 {
		log.Fatal("all workers returned error on startup")
	}

	// proceeed Tests with workers that started correctly
	common.Init(startedWorkers)
	os.Exit(m.Run())
}
