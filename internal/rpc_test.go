package internal

import (
	"distributed/common"
	"testing"
)

func TestRpcNode(t *testing.T) {

	worker := common.GetAvailRegWorkers()[0]
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
	nodesToDel := []RpcNode{nodeInfo}
	err = worker.Invoke(DELETE, nodesToDel, &common.NONE{})
	if err != nil {
		t.Fail()
	}
	err = worker.Invoke(RETRIEVE, nodeInfo, &result)
	if err == nil || err == common.NoResultsErr {
		t.Fail()
	}

}
