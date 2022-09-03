/*	@author: Bilal El Uneis
	@since: July 2022
	@email: bilaleluneis@gmail.com	*/

package datastructure

import (
	"distributed/common"
	"distributed/internal"
)

// UpdateNode will update existing node provided GRP ID and UUID
// via finding, deleting then inserting again
// TODO: no very performant impl as I am making 3 calls on all workers
func UpdateNode[T any](withNode common.NodeLike[T]) error {
	var err error
	// find first if node with GRP ID and UUID exist
	if _, err = FindNodeByUuid[T](withNode.GetUuID(), withNode.GetGrpID()); err == nil {
		if err = DeleteNodes([]common.UUID{withNode.GetUuID()}, withNode.GetGrpID()); err == nil {
			if updatedNode := internal.Encode[T](withNode); len(updatedNode) != 0 {
				var worker common.RegisteredWorker
				if worker, err = common.GetRandomAvailRegWorker(); err == nil {
					err = worker.Invoke(internal.INSERT, updatedNode[0], &common.NONE{})
				}
			}
		}
	}
	return err
}

func DeleteNodes(uuids []common.UUID, forGrp common.GRPID) error {
	if forGrp == common.EmptyGrpID {
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
func FindNodeByUuid[T any](uuid common.UUID, forGrp common.GRPID) (common.Node[T], error) {
	searchParam := internal.RpcNode{
		GrpID: forGrp,
		Uuid:  uuid,
	}
	nodesFound := retrieveFromWorkers[T](searchParam)
	if len(nodesFound) == 0 {
		return common.Node[T]{}, common.NoResultsErr
	}
	if len(nodesFound) > 1 {
		return common.Node[T]{}, common.MultipleMatchErr
	}
	return nodesFound[0], nil
}

func FindNodesByValue[T any](value T, forGroup common.GRPID) ([]common.Node[T], error) {
	var err error
	var bytesVal []byte
	result := make([]common.Node[T], 0)
	if bytesVal, err = common.ToBytes[T](value); err != nil {
		return result, err
	}
	searchParam := internal.RpcNode{
		GrpID: forGroup,
		Data:  bytesVal,
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
	var worker common.RegisteredWorker

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
	node := common.Node[T]{
		Data:  withVal,
		GrpId: grpId,
		Uuid:  uuid,
	}
	if rpcNodes := internal.Encode[T](node); len(rpcNodes) != 0 {
		if worker, err = common.GetRandomAvailRegWorker(); err != nil {
			return common.EmptyGrpID, common.EmptyUUID, err
		}
		if err = worker.Invoke(internal.INSERT, rpcNodes[0], &common.NONE{}); err != nil {
			return common.EmptyGrpID, common.EmptyUUID, err
		}
	}
	return grpId, uuid, err
}

// calls Retrieve internal.rpc on all workers and return slice of found results
// workers that return error on call will be skipped for result
// TODO: consider using channels and goroutines/ async calls
func retrieveFromWorkers[T any](searchParm internal.RpcNode) []common.Node[T] {
	var result []common.Node[T]
	workers := common.GetAvailRegWorkers()
	nodesFound := make([]internal.RpcNode, 0)
	for _, worker := range workers {
		if err := worker.Invoke(internal.RETRIEVE, searchParm, &nodesFound); err == nil {
			result = append(result, internal.Decode[T](nodesFound...)...)
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
			return common.EmptyGrpID, err
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
