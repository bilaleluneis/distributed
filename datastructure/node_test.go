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

func TestTcpRpc(t *testing.T) {

	node := Node[int]{
		IsRoot: true,
		Data:   1,
	}

	location := common.Location{
		HostName: "localhost",
		Ip:       8080,
	}

	// test insert op
	uuid, err := Insert(node, location)
	if err != nil || uuid == "" {
		t.Fail()
	}

	// test retrieve op
	location.Uuid = uuid
	params := common.SearchParams{
		Remove:  false,
		Address: location,
	}
	returnedNode, err := Retrieve[int](params)
	if err != nil || returnedNode.Data != 1 {
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

	exitVal := m.Run()
	os.Exit(exitVal)
}
