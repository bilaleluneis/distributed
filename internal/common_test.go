/*	@author: Bilal El Uneis
	@since: Sept 2022
	@email: bilaleluneis@gmail.com	*/

package internal

import (
	"distributed/common"
)

type vegetableFilter struct{}

func (fi vegetableFilter) Filter(item common.NodeLike[string]) bool {
	return item.GetData() == "lettuce"
}

// Utils

func genTestNodes[T any](data ...T) (common.GRPID, error) {
	var err error
	var worker common.RegisteredWorker
	var grpId common.GRPID

	if worker, err = common.GetRandomAvailRegWorker(); err != nil {
		return common.EmptyGrpID, err
	}

	// convert data passed to []byte to be used in RpcNode type
	testDatas := make([][]byte, 0)
	for _, item := range data {
		var b []byte
		if b, err = common.ToBytes(item); err == nil {
			testDatas = append(testDatas, b)
		}
	}

	// create nodes with Test data
	if err = worker.Invoke(NEW, common.NONE{}, &grpId); err != nil {
		return common.EmptyGrpID, err
	}
	for _, testdata := range testDatas {
		node := RpcNode{
			Data:  testdata,
			GrpID: grpId,
			Uuid:  common.GenUUID(),
		}
		if err = worker.Invoke(INSERT, node, &common.NONE{}); err != nil {
			return common.EmptyGrpID, err
		}
	}

	return grpId, err
}
