package internal

import (
	"distributed/common"
	"log"
	"os"
	"testing"
)

func TestRpcNode(t *testing.T) {

	worker := common.Worker{
		Host: "localhost",
		Port: 8080,
	}

	var grpId common.GRPID
	err := worker.Invoke(NEW, common.NONE{}, &grpId)

	node := RpcNode{
		Data:  []byte("testData"),
		GrpID: grpId,
		Uuid:  common.GenUUID(),
	}

	// Test Insert
	err = worker.Invoke(INSERT, node, &common.NONE{})
	if err != nil {
		t.Fail()
	}

	nodeInfo := RpcNode{
		GrpID: grpId,
		Uuid:  node.Uuid,
	}
	var result []RpcNode

	// Test Retrieve
	err = worker.Invoke(RETRIEVE, nodeInfo, &result)
	if err != nil {
		t.Fail()
	}
	if len(result) == 0 {
		t.Fail()
	}

	// Test Delete
	err = worker.Invoke(DELETE, nodeInfo, &common.NONE{})
	if err != nil {
		t.Fail()
	}
	err = worker.Invoke(RETRIEVE, nodeInfo, &result)
	if err == nil || err == common.NoResultsErr {
		t.Fail()
	}

}

func TestMain(m *testing.M) {

	go func() {
		err := InitServer(8080)
		if err != nil {
			log.Fatal("test server startup error", err)
		}
	}()

	os.Exit(m.Run())

}
