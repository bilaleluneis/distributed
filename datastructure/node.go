/*	@author: Bilal El Uneis
	@since: July 2022
	@email: bilaleluneis@gmail.com	*/

package datastructure

import (
	"bytes"
	"distributed/common"
	"distributed/internal"
	"encoding/gob"
	"fmt"
)

//TODO: refactor to use NodeID type

type Node[T any] struct {
	Data   T
	GrpId  common.GRPID
	Uuid   common.UUID
	Parent common.UUID
	Child  common.UUID
}

func (n Node[T]) String() string {
	rep := "Node[T]: grpId=%s uuid=%s parent=%s child=%s"
	return fmt.Sprintf(rep, n.GrpId, n.Uuid, n.Parent, n.Child)
}

// UpdateNode will update existing node provided GRP ID and UUID
// TODO: no very performant impl as I am making 3 calls on all workers
// find, delete and insert
func UpdateNode[T any](withNode Node[T]) error {
	var err error
	// find first if node with GRP ID and UUID exist
	if _, err = FindNodeByUuid[T](withNode.Uuid, withNode.GrpId); err == nil {
		if err = DeleteNodes([]common.UUID{withNode.Uuid}, withNode.GrpId); err == nil {
			if updatedNode := encode[T]([]Node[T]{withNode}); len(updatedNode) != 0 {
				err = common.GetRandomAvailRegWorker().Invoke(internal.INSERT, updatedNode[0], &common.NONE{})
			}
		}
	}
	return err
}

func DeleteNodes(uuids []common.UUID, forGrp common.GRPID) error {
	if forGrp == "" {
		return common.ReqGrpIdErr
	}
	if len(uuids) == 0 {
		return common.NonToDelErr
	}
	nodesToDel := make([]internal.RpcNode, 0)
	for _, uuid := range uuids {
		nodesToDel = append(nodesToDel, internal.RpcNode{GrpID: forGrp, Uuid: uuid})
	}
	workers := common.GetAvailRegWorkers()
	for _, worker := range workers {
		_ = worker.Invoke(internal.DELETE, nodesToDel, &common.NONE{})
	}
	return nil
}

// FindNodeByUuid assumption here is that there should be one node only with that match uuid
func FindNodeByUuid[T any](uuid common.UUID, forGrp common.GRPID) (Node[T], error) {
	searchParam := internal.RpcNode{
		GrpID: forGrp,
		Uuid:  uuid,
	}
	nodesFound := retrieveFromWorkers[T](searchParam)
	if len(nodesFound) == 0 {
		return Node[T]{}, common.NoResultsErr
	}
	if len(nodesFound) > 1 {
		return Node[T]{}, common.MultipleMatchErr
	}
	return nodesFound[0], nil
}

func FindNodesByValue[T any](data T, forGroup common.GRPID) ([]Node[T], error) {
	var result []Node[T]
	var err error
	buffer := bytes.NewBuffer([]byte{})
	err = gob.NewEncoder(buffer).Encode(data)
	if err != nil {
		return make([]Node[T], 0), err
	}
	searchParam := internal.RpcNode{
		GrpID: forGroup,
		Data:  buffer.Bytes(),
	}
	if result = retrieveFromWorkers[T](searchParam); len(result) == 0 {
		err = common.NoResultsErr
	}
	return result, err
}

func NewNode[T any](withVal T, inGrp common.GRPID) (common.GRPID, common.UUID, error) {
	grpId := inGrp
	var uuid common.UUID
	var err error

	workers := common.GetAvailRegWorkers()
	if len(workers) == 0 {
		return grpId, uuid, common.NoWorkerAvailErr
	}
	if grpId == common.EmptyGrpID {
		if grpId, err = genGroupID(workers); err != nil {
			return grpId, uuid, err
		}
	}
	if uuid, err = genUUID(workers, grpId); err != nil {
		return grpId, uuid, err
	}
	// persist node
	node := Node[T]{
		Data:  withVal,
		GrpId: grpId,
		Uuid:  uuid,
	}
	if rpcNodes := encode[T]([]Node[T]{node}); len(rpcNodes) != 0 {
		worker := common.GetRandomAvailRegWorker()
		rpcNode := rpcNodes[0]
		err = worker.Invoke(internal.INSERT, rpcNode, &common.NONE{})
	}

	return grpId, uuid, err
}

// calls Retrieve internal.rpc on all workers and return slice of found results
// workers that return error on call will be skipped for result
// TODO: consider using channels and goroutines/ async calls
func retrieveFromWorkers[T any](searchParm internal.RpcNode) []Node[T] {
	result := make([]Node[T], 0)
	workers := common.GetAvailRegWorkers()
	for _, worker := range workers {
		nodesFound := make([]internal.RpcNode, 0)
		err := worker.Invoke(internal.RETRIEVE, searchParm, &nodesFound)
		if err == nil {
			for _, n := range nodesFound {
				if ndec := decode[T]([]internal.RpcNode{n}); len(ndec) != 0 {
					result = append(result, ndec[0])
				}
			}
		}
	}
	return result
}

func genGroupID(workers []common.RegisteredWorker) (common.GRPID, error) {
genID:
	grpId := common.GenUUID()
	var err error
	for _, worker := range workers {
		var exist bool
		err = worker.Invoke(internal.GRPIDEXIST, grpId, &exist)
		if err != nil {
			return "", err
		}
		if exist {
			goto genID
		}
	}
	return grpId, err
}

func genUUID(workers []common.RegisteredWorker, forGrp common.GRPID) (common.UUID, error) {
genID:
	uuid := common.GenUUID()
	node := internal.RpcNode{
		GrpID: forGrp,
		Uuid:  uuid,
	}
	for _, worker := range workers {
		var nodesFound []internal.RpcNode
		err := worker.Invoke(internal.RETRIEVE, node, &nodesFound)
		if err == nil && len(nodesFound) != 0 {
			goto genID
		}
	}
	return uuid, nil
}

func encode[T any](nodes []Node[T]) []internal.RpcNode {
	result := make([]internal.RpcNode, 0)
	for _, n := range nodes {
		var buffer bytes.Buffer
		err := gob.NewEncoder(&buffer).Encode(n.Data)
		if err == nil {
			result = append(result, internal.RpcNode{
				Data:   buffer.Bytes(),
				GrpID:  n.GrpId,
				Uuid:   n.Uuid,
				Parent: n.Parent,
				Child:  n.Child,
			})
		}
	}
	return result
}

func decode[T any](rpcNodes []internal.RpcNode) []Node[T] {
	result := make([]Node[T], 0)
	for _, n := range rpcNodes {
		var node Node[T]
		buffer := bytes.NewBuffer(n.Data)
		err := gob.NewDecoder(buffer).Decode(&node.Data)
		if err == nil {
			node.GrpId = n.GrpID
			node.Uuid = n.Uuid
			node.Parent = n.Parent
			node.Child = n.Child
			result = append(result, node)
		}
	}
	return result
}
