/*	@author: Bilal El Uneis
	@since: Sep 2022
	@email: bilaleluneis@gmail.com	*/

package datastructure

import (
	"distributed/common"
	"distributed/internal"
	"encoding/gob"
)

func Filter[T any](c common.Collection, f common.Filterer[T]) error {
	filter := internal.Filter[T]{WithFilter: f}
	gob.Register(f)
	gob.Register(filter)
	param := internal.FuncParam{
		GrpId: c.Identity(),
		Op:    filter,
	}
	workers := common.GetAvailRegWorkers()
	for _, worker := range workers {
		if err := worker.Invoke(internal.DELAYED, &param, &common.NONE{}); err != nil {
			return FILTEROPERR //TODO use purge on workers to clean up
		}
	}
	return nil
}

func Map[T any, R any](c common.Collection, m common.Mapper[T, R]) error {
	maper := internal.Map[T, R]{WithMapper: m}
	gob.Register(m)
	gob.Register(maper)
	param := internal.FuncParam{
		GrpId: c.Identity(),
		Op:    maper,
	}
	workers := common.GetAvailRegWorkers()
	for _, worker := range workers {
		if err := worker.Invoke(internal.DELAYED, &param, &common.NONE{}); err != nil {
			return MAPOPERR //TODO use purge on workers to clean up
		}
	}
	return nil
}

func Reduce[T any, R any](c common.Collection, r common.Reducer[T, R], finalReduct func([]R) R) (R, error) {
	var err error
	reduce := internal.Reduce[T, R]{WithReducer: r}
	gob.Register(r)
	gob.Register(reduce)
	param := internal.FuncParam{
		Op:    reduce,
		GrpId: c.Identity(),
	}
	workers := common.GetAvailRegWorkers()

	// capture reduction result from each worker and place into slice
	workersRedResult := make([]internal.RpcNode, 0)
	for _, worker := range workers {
		var currWorkerReduct []byte
		if err = worker.Invoke(internal.IMMEDIATE, &param, &currWorkerReduct); err == nil {
			var toRpcNodes []internal.RpcNode
			if toRpcNodes, err = common.ToType[[]internal.RpcNode](currWorkerReduct); err == nil {
				workersRedResult = append(workersRedResult, toRpcNodes...)
			}
		}
	}

	var result R
	if err != nil {
		return result, err
	}

	switch len(workersRedResult) {
	case 0:
		return result, REDUCEOPERR
	case 1:
		result, err = common.ToType[R](workersRedResult[0].Data)
	default:
		intermReduction := make([]R, 0)
		for _, cr := range workersRedResult {
			var r R
			r, err = common.ToType[R](cr.Data)
			if err != nil {
				return result, err
			}
			intermReduction = append(intermReduction, r)
		}
		result = finalReduct(intermReduction)
	}

	return result, err
}

func Compute[T any](c common.Collection) ([]T, error) {
	var err error
	compute := internal.Compute{}
	gob.Register(compute)
	param := internal.FuncParam{
		GrpId: c.Identity(),
		Op:    compute,
	}

	workers := common.GetAvailRegWorkers()
	computeResult := make([]internal.RpcNode, 0)
	for _, worker := range workers {
		var tmp []byte
		if err = worker.Invoke(internal.IMMEDIATE, &param, &tmp); err == nil {
			var result []internal.RpcNode
			if result, err = common.ToType[[]internal.RpcNode](tmp); err == nil {
				computeResult = append(computeResult, result...)
			}
		}
	}

	if err != nil {
		return []T{}, COMPUTEOPERR
	}

	result := make([]T, 0)
	if len(computeResult) != 0 {
		for _, cr := range computeResult {
			var r T
			if r, err = common.ToType[T](cr.Data); err == nil {
				result = append(result, r)
			}
		}
	}

	if err != nil {
		return []T{}, COMPUTEOPERR
	}

	return result, nil
}
