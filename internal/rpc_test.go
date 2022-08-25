/*	@author: Bilal El Uneis
	@since: June 2022
	@email: bilaleluneis@gmail.com	*/

package internal

import (
	"bytes"
	"distributed/common"
	"encoding/gob"
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

// Filterer Impl to use
type vegetableFilter struct{}

func (fi vegetableFilter) Func(nodes []RpcNode) []RpcNode {
	result := make([]RpcNode, 0)
	for _, node := range nodes {
		if bytes.Equal(node.Data, []byte("lettuce")) {
			result = append(result, node)
		}
	}
	return result
}

type doNothingReduce struct{}

func (r doNothingReduce) Func(nodes []RpcNode) []RpcNode {
	return nodes
}

// To understand how interfaces work with GOB
// look at research/gob_test.go TestInterfaceOverGob
func TestFilter(t *testing.T) {
	var err error
	var grpId common.GRPID
	worker := common.GetAvailRegWorkers()[0]

	// Data setup
	if err = worker.Invoke(NEW, common.NONE{}, &grpId); err == nil {
		testDatas := [][]byte{
			[]byte("apple"),
			[]byte("orange"),
			[]byte("lemon"),
			[]byte("peach"),
			[]byte("lettuce"),
		}
		for _, testdata := range testDatas {
			node := RpcNode{
				Data:  testdata,
				GrpID: grpId,
				Uuid:  common.GenUUID(),
			}
			if err = worker.Invoke(INSERT, node, &common.NONE{}); err != nil {
				t.FailNow()
			}
		}
	}

	// Test starts here

	// register types for Functional Interface
	gob.Register(vegetableFilter{})
	gob.Register(doNothingReduce{})

	param := FunctionParam{
		GrpID:    grpId,
		Function: &vegetableFilter{},
	}
	// must pass ref to interface to work with GOB
	if err = worker.Invoke(FILTER, &param, &common.NONE{}); err == nil {
		param.Function = &doNothingReduce{}
		var result []RpcNode
		if err = worker.Invoke(REDUCE, &param, &result); err == nil {
			if len(result) != 1 || !bytes.Equal(result[0].Data, []byte("lettuce")) {
				t.Fail()
			}
		}
		return // got here Test succeeded
	}

	// got here , worker invokation failed
	t.Fatalf("worker invoke failure %s", err.Error()) // got here
}
