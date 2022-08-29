/*	@author: Bilal El Uneis
	@since: June 2022
	@email: bilaleluneis@gmail.com	*/

package datastructure

import (
	"distributed/common"
	"log"
	"testing"
)

func TestNewNode(t *testing.T) {
	if grpId, uuid, err := NewNode[int](1, common.EmptyGrpID); err == nil {
		if grpId == common.EmptyGrpID {
			log.Printf("empty group id")
			t.Fail()
		}
		if uuid == common.EmptyUUID {
			log.Printf("empty uuid")
			t.Fail()
		}
	} else {
		t.Fatalf("NewNode failed with error %s", err.Error())
	}
}

func TestFindNodeByUuid(t *testing.T) {
	var grpId common.GRPID
	var uuid common.UUID
	var err error
	var node Node[int]
	value := 1

	// NOTE: if failed here then check TestNewNode
	if grpId, uuid, err = NewNode[int](value, common.EmptyGrpID); err != nil {
		t.Fatalf("creating new node failed due %s", err.Error())
	}

	// Test find node by uuid starts
	if node, err = FindNodeByUuid[int](uuid, grpId); err != nil {
		t.Fatalf("finding node by uuid with error %s", err.Error())
	}
	if node.Uuid != uuid || node.GrpId != grpId {
		t.Fatal("node uuid and/or node grpId do not match request")
	}
	if node.Data != value {
		t.Fatalf("value of node doesnt match expected %d got %d", value, node.Data)
	}
}

func TestFindNodesByValue(t *testing.T) {
	var grpId common.GRPID
	var err error
	var nodes []Node[int]
	values := []int{2, 3, 4, 2, 6, 7, 8, 9, 2}

	// Setup Test Data
	for _, value := range values {
		if grpId, _, err = NewNode[int](value, grpId); err != nil {
			t.Fatalf("creating new nodes failed due %s", err.Error())
		}
	}

	// Test find nodes by value
	nodes, err = FindNodesByValue[int](2, grpId)
	if err != nil {
		t.Fatalf("finding node by value with error %s", err.Error())
	}
	if len(nodes) != 3 {
		t.Fatalf("expected 3 nodes got %d", len(nodes))
	}
	for _, node := range nodes {
		if node.Data != 2 {
			t.Fatalf("%s has value of %d expected 2", node, node.Data)
		}
	}
}

func TestUpdateNode(t *testing.T) {
	var grpId common.GRPID
	var uuid common.UUID
	var err error

	// Setup Test Data
	if grpId, uuid, err = NewNode[int](1, grpId); err != nil {
		t.Fatalf("creating new nodes failed due %s", err.Error())
	}

	// Test Update
	err = UpdateNode[int](Node[int]{
		GrpId: grpId,
		Uuid:  uuid,
		Data:  2,
	})
	if err != nil {
		t.Fatalf("update node failed due %s", err.Error())
	}

	// Retrieve and verify update
	updatedNode, err := FindNodeByUuid[int](uuid, grpId)
	if err != nil {
		t.Fatalf("retrieve updated node failed due %s", err.Error())
	}
	if updatedNode.Data != 2 {
		t.Fatalf("updat node failed expected value of 2 got %d", updatedNode.Data)
	}
}

func TestDeleteNodes(t *testing.T) {
	var grpId common.GRPID
	var err error
	var nodes []Node[int]
	values := []int{2, 3, 4, 2, 6, 7, 8, 9, 2}

	// Setup Test Data
	for _, value := range values {
		if grpId, _, err = NewNode[int](value, grpId); err != nil {
			t.Fatalf("creating new nodes failed due %s", err.Error())
		}
	}

	// Find nodes to delete
	nodes, _ = FindNodesByValue[int](2, grpId)
	uuids := make([]common.UUID, 0)
	for _, node := range nodes {
		uuids = append(uuids, node.Uuid)
	}

	// Test Delete
	err = DeleteNodes(uuids, grpId)
	if err != nil {
		t.Fatalf("delete nodes failed due %s", err.Error())
	}
	if _, err := FindNodesByValue(2, grpId); err != common.NoResultsErr {
		t.Fatal("expected no result error")
	}
}
