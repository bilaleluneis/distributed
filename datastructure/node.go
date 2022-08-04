/*	@author: Bilal El Uneis
	@since: July 2022
	@email: bilaleluneis@gmail.com	*/

package datastructure

import (
	"bytes"
	"distributed/common"
	"distributed/internal"
	"encoding/gob"
)

type Node[T any] struct {
	Data   T
	GrpId  common.GRPID
	Uuid   common.UUID
	Parent common.UUID
	Child  common.UUID
}

func New[T any](data T, inGrp common.GRPID) (common.GRPID, common.UUID, error) {
	grpId := inGrp
	var uuid common.UUID
	var err error

	workers := common.Get()
	if len(workers) == 0 {
		return grpId, uuid, common.NoWorkerAvailErr
	}
	if grpId == "" {
		grpId, err = genGroupID(workers)
		if err != nil {
			return grpId, uuid, err
		}
	}

	uuid, err = genUUID(workers, grpId)
	if err != nil {
		return grpId, uuid, err
	}

	// persist node
	buffer := bytes.NewBuffer([]byte{})
	err = gob.NewEncoder(buffer).Encode(data)
	if err != nil {
		return grpId, uuid, err
	}
	node := internal.RpcNode{
		GrpID: grpId,
		Uuid:  uuid,
		Data:  buffer.Bytes(),
	}
	err = workers[0].Invoke(internal.INSERT, node, &common.NONE{})

	return grpId, uuid, err
}

func genGroupID(workers []common.Worker) (common.GRPID, error) {
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

func genUUID(workers []common.Worker, forGrp common.GRPID) (common.UUID, error) {
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
