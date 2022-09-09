/*	@author: Bilal El Uneis
	@since: Sep 2022
	@email: bilaleluneis@gmail.com	*/

package datastructure

import (
	"distributed/common"
	"distributed/internal"
)

func Filter[T any](c common.Collection, f common.Filterer[T]) error {
	filter := internal.Filter[T]{WithFilter: f}
	if err := delayedEval(filter, c.Identity()); err != nil {
		return FILTEROPERR
	}
	return nil
}

func Map[T any, R any](c common.Collection, m common.Mapper[T, R]) error {
	mapper := internal.Map[T, R]{WithMapper: m}
	if err := delayedEval(mapper, c.Identity()); err != nil {
		return MAPOPERR
	}
	return nil
}

func Reduce[T any, R any](c common.Collection, r common.Reducer[T, R], finalReduct func([]R) R) (R, error) {
	var err error
	var result R
	var workersResult []internal.RpcNode
	reduce := internal.Reduce[T, R]{WithReducer: r}

	if workersResult, err = eagerEval(reduce, c.Identity()); err != nil {
		return result, err
	}

	switch len(workersResult) {
	case 0:
		return result, REDUCEOPERR
	case 1:
		result, err = common.ToType[R](workersResult[0].Data)
	default:
		intermReduction := make([]R, 0)
		for _, cr := range workersResult {
			var r R
			if r, err = common.ToType[R](cr.Data); err == nil {
				intermReduction = append(intermReduction, r)
			} else {
				return result, err
			}
		}
		result = finalReduct(intermReduction)
	}

	return result, err
}

func Compute[T any](c common.Collection) ([]T, error) {
	var err error
	var computeResult []internal.RpcNode
	compute := internal.Compute{}

	if computeResult, err = eagerEval(compute, c.Identity()); err != nil {
		return []T{}, err
	}

	result := make([]T, 0)
	if len(computeResult) != 0 {
		for _, cr := range computeResult {
			var r T
			if r, err = common.ToType[T](cr.Data); err == nil {
				result = append(result, r)
			} else {
				return []T{}, err
			}
		}
	}

	return result, nil
}

func eagerEval(fop internal.FunctionalOp, forGrp common.GRPID) ([]internal.RpcNode, error) {
	param := internal.FuncParam{
		Op:    fop,
		GrpId: forGrp,
	}
	workers := common.GetAvailRegWorkers()
	result := make([]internal.RpcNode, 0)
	for _, worker := range workers {
		var currWorkerReduct []internal.RpcNode
		if err := worker.Invoke(internal.IMMEDIATE, &param, &currWorkerReduct); err == nil {
			result = append(result, currWorkerReduct...)
		} else {
			return []internal.RpcNode{}, err
		}
	}
	return result, nil
}

func delayedEval(fop internal.FunctionalOp, forGrp common.GRPID) error {
	param := internal.FuncParam{
		GrpId: forGrp,
		Op:    fop,
	}
	workers := common.GetAvailRegWorkers()
	for _, worker := range workers {
		if err := worker.Invoke(internal.DELAYED, &param, &common.NONE{}); err != nil {
			return err
		}
	}
	return nil
}
