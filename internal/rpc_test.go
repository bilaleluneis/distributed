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
func TestFilterReduceOp(t *testing.T) {
	var err error
	var worker common.RegisteredWorker
	var grpId common.GRPID
	fruitesAndVeg := []string{
		"apple",
		"orange",
		"lemon",
		"peach",
		"lettuce"}

	// gen test data, if you get failure here make sure
	// TestCreateRetrieve is passing before debugging futher
	grpId, err = genTestNodes(fruitesAndVeg...)

	// Test Filter
	filter := Filter[vegetableFilter, string]{}
	gob.Register(filter)
	param := FuncParam{
		Op:    &filter,
		GrpId: grpId,
	}
	if worker, err = common.GetRandomAvailRegWorker(); err != nil {
		t.Fatalf("obtain worker %s", err.Error())
	}
	if err = worker.Invoke(FILTER, &param, &common.NONE{}); err != nil {
		t.Fatalf("Filter RPC call %s", err.Error())
	}

	// reduce to get result
	reduce := Reduce[countReducer, string, int]{}
	gob.Register(reduce)
	param = FuncParam{
		Op:    &reduce,
		GrpId: grpId,
	}
	var rpcCallResult []byte
	if err = worker.Invoke(REDUCE, &param, &rpcCallResult); err != nil {
		t.Fatalf("error %s", err.Error())
	}
	count, err := common.ToType[int](rpcCallResult)
	if err != nil {
		t.Fatalf("conversion to int %s", err.Error())
	}
	if count != 1 {
		t.Fatalf("count expected %d got %d", 1, count)
	}
}
