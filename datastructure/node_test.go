/*	@author: Bilal El Uneis
	@since: June 2022
	@email: bilaleluneis@gmail.com	*/

package datastructure

import (
	"distributed/common"
	"testing"
)

// TODO: clean up test with better details and use cases
// add tests with multiple workers
func TestNode(t *testing.T) {
	var err error
	// Test creating a node
	grpId, uuid, err := NewNode(1, "")
	if err != nil {
		t.Fatalf("Create Node: %s", err)
	}
	if grpId == "" || uuid == "" {
		t.Fatalf("Create Node GrpID or UUID is empty")
	}

	// Test find node by value
	nodesFound, err := FindNodesByValue(1, grpId)
	if err != nil {
		t.Fatal(err)
	}
	if len(nodesFound) != 1 {
		t.Fatalf("multiple results")
	}
	if nodesFound[0].Data != 1 {
		t.Fatalf("incorrect value for Data")
	}

	//Test find node by uuid
	nodeFound, err := FindNodeByUuid[int](nodesFound[0].Uuid, nodesFound[0].GrpId)
	if err != nil {
		t.Fatal(err)
	}
	if nodeFound.Uuid != nodesFound[0].Uuid {
		t.Fatalf("incorrect value for Data")
	}

	// Test update node
	nodeFound.Data = 2
	err = UpdateNode(nodeFound)
	if err != nil {
		t.Fatal(err)
	}
	nodeFound, err = FindNodeByUuid[int](nodeFound.Uuid, nodeFound.GrpId)
	if err != nil {
		t.Fatal(err)
	}
	if nodeFound.Data != 2 {
		t.Fatalf("node was not updated")
	}

	//Test Delete node
	_ = DeleteNodes([]common.UUID{nodeFound.Uuid}, nodeFound.GrpId)
	_, err = FindNodeByUuid[int](nodeFound.Uuid, nodeFound.GrpId)
	if err != common.NoResultsErr {
		t.Fatalf("found a Deleted node")
	}
}
