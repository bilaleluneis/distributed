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

	node := Node[int]{
		Data: 1,
	}

	location := common.Location{
		HostName: "localhost",
		Port:     8080,
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

	//test update op
	returnedNode.Data = 2
	err = Update[int](location, returnedNode)
	if err != nil {
		t.Fail()
	}
	updatedNode, err := Retrieve[int](params)
	if err != nil || updatedNode.Data != 2 {
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
