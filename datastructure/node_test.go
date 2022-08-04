/*	@author: Bilal El Uneis
	@since: June 2022
	@email: bilaleluneis@gmail.com	*/

package datastructure

import (
	"distributed/common"
	"distributed/internal"
	"log"
	"os"
	"testing"
)

func TestNode(t *testing.T) {
	grpId, uuid, err := New[int](1, "")
	if err != nil {
		t.Fail()
	}
	if grpId == "" || uuid == "" {
		t.Fail()
	}
}

// allow single initialization of client and server instance for all tests
func TestMain(m *testing.M) {

	go func() {
		err := internal.InitServer(8080)
		if err != nil {
			log.Fatal("test server startup error", err)
		}
	}()

	workers := []common.Worker{{
		Host: "localhost",
		Port: 8080,
	}}
	common.Init(workers)

	os.Exit(m.Run())
}
