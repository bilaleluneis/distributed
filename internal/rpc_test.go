/*	@author: Bilal El Uneis
	@since: June 2022
	@email: bilaleluneis@gmail.com	*/

package internal

import (
	"distributed/common"
	"encoding/gob"
	"testing"
)

func TestCreateRetrieve(t *testing.T) {
	var err error
	var worker common.RegisteredWorker
	var grpId common.GRPID
	var sentData []byte
	var returnedData int
	var retResult []RpcNode

	if worker, err = common.GetRandomAvailRegWorker(); err != nil {
		t.Fatal(err)
	}
	if err = worker.Invoke(NEW, common.NONE{}, &grpId); err != nil {
		t.Fatalf("NEW with error %s", err.Error())
	}

	//setup Data to insert
	sentData, err = common.ToBytes[int](1)
	node := RpcNode{
		Data:  sentData,
		GrpID: grpId,
		Uuid:  common.GenUUID(),
	}

	// Test Insert
	if err = worker.Invoke(INSERT, node, &common.NONE{}); err != nil {
		t.Fatalf("INSERT with error %s", err.Error())
	}

	// validate via Retrieve
	nodeInfo := RpcNode{
		GrpID: grpId,
		Uuid:  node.Uuid,
	}

	// Test Retrieve
	if err = worker.Invoke(RETRIEVE, nodeInfo, &retResult); err != nil {
		t.Fatalf("RETRIEVE with error %s", err.Error())
	}
	returnedData, err = common.ToType[int](retResult[0].Data)
	if err != nil {
		t.Fatalf("conversion back failed with error %s", err.Error())
	}
	if returnedData != 1 {
		t.Fatalf("conversion back failed expected 1 got %d", returnedData)
	}
}

// To understand how interfaces work with GOB
// look at research/gob_test.go TestInterfaceOverGob
func TestFilter(t *testing.T) {
	var err error
	var worker common.RegisteredWorker
	var grpId common.GRPID
	var rpcNodes []RpcNode
	fruitesAndVeg := []string{
		"apple",
		"orange",
		"lemon",
		"peach",
		"lettuce"}

	if worker, err = common.GetRandomAvailRegWorker(); err != nil {
		t.Fatalf("obtain worker %s", err.Error())
	}

	// gen test data, if you get failure here make sure
	// TestCreateRetrieve is passing before debugging futher
	grpId, err = genTestNodes(fruitesAndVeg...)

	// Test Filter
	filter := Filter[string]{vegetableFilter{}}
	gob.Register(vegetableFilter{})
	gob.Register(filter)
	param := FuncParam{
		Op:    &filter,
		GrpId: grpId,
	}
	var byteResult []byte
	if err = worker.Invoke(IMMEDIATE, &param, &byteResult); err != nil {
		t.Fatalf("filter call %s", err.Error())
	}
	if rpcNodes, err = common.ToType[[]RpcNode](byteResult); err != nil {
		t.Fatalf("convertion from []byts to []RpcNode %s", err.Error())
	}
	nodes := Decode[string](rpcNodes...)
	if len(nodes) != 1 {
		t.Fatalf("wrong number of result")
	}
	if result := nodes[0].Data; result != "lettuce" {
		t.Fatalf("wrong filter size or data")
	}
}
