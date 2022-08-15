/*	@author: Bilal El Uneis
	@since: August 2022
	@email: bilaleluneis@gmail.com	*/

package common

import "testing"

func TestWorkerRegistration(t *testing.T) {

	availableWorkers := GetAvailRegWorkers()
	numAvailWorkers := len(availableWorkers)

	// Test no available workers
	if numAvailWorkers != 0 {
		t.Fatalf("expected initial available workers to be zero got %d", numAvailWorkers)
	}

	// Test registering workers
	workers := []struct {
		host string
		port int
	}{
		{"localhost", 8080},
		{"localhost", 8081},
		{"localhost", 8082},
	}
	for _, worker := range workers {
		RegisterWorker(worker.host, worker.port)
	}
	availableWorkers = GetAvailRegWorkers()
	numAvailWorkers = len(availableWorkers)
	if numAvailWorkers != len(workers) {
		t.Fatalf("expected %d workers but got %d", len(workers), numAvailWorkers)
	}

	//TODO:
	// Test Random Worker

}
