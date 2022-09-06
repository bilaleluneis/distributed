/*	@author: Bilal El Uneis
	@since: August 2022
	@email: bilaleluneis@gmail.com	*/

package internal

import (
	"distributed/common"
)

type FunctionalOp interface {
	Eval([]RpcNode) []RpcNode
}

type FuncParam struct {
	Op    FunctionalOp
	GrpId common.GRPID
}

type Compute struct{}

func (Compute) Eval(rpcNodes []RpcNode) []RpcNode {
	return rpcNodes
}

type Filter[T any] struct {
	WithFilter common.Filterer[T]
}

func (f Filter[T]) Eval(rpcNodes []RpcNode) []RpcNode {
	result := make([]RpcNode, 0)
	for _, rpcNode := range rpcNodes {
		n := Decode[T](rpcNode)
		if f.WithFilter.Filter(n[0]) {
			result = append(result, rpcNode)
		}
	}
	return result
}

type Reduce[T any, R any] struct {
	WithReducer common.Reducer[T, R]
}

func (r Reduce[T, R]) Eval(rpcNodes []RpcNode) []RpcNode {
	nodes := Decode[T](rpcNodes...)
	reductionList := make([]common.NodeLike[T], 0)
	for _, n := range nodes {
		reductionList = append(reductionList, n)
	}
	return []RpcNode{Encode[R](common.Node[R]{
		Data: r.WithReducer.Reduce(reductionList...),
	})[0]}
}

type Map[T any, R any] struct {
	WithMapper common.Mapper[T, R]
}

func (m Map[T, R]) Eval(rpcNodes []RpcNode) []RpcNode {
	result := make([]RpcNode, 0)
	for _, rpcNode := range rpcNodes {
		node := Decode[T](rpcNode)[0]
		mappedNode := m.WithMapper.Map(&node)
		result = append(result, Encode[R](mappedNode)[0])
	}
	return result
}

func Encode[T any](nodes ...common.NodeLike[T]) []RpcNode {
	result := make([]RpcNode, 0)
	for _, n := range nodes {
		if data, err := common.ToBytes[T](n.GetData()); err == nil {
			result = append(result, RpcNode{
				Data:   data,
				GrpID:  n.GetGrpID(),
				Uuid:   n.GetUuID(),
				Parent: n.GetParent(),
				Child:  n.GetChild(),
			})
		}
	}
	return result
}

func Decode[T any](rpcNodes ...RpcNode) []common.Node[T] {
	result := make([]common.Node[T], 0)
	for _, n := range rpcNodes {
		if value, err := common.ToType[T](n.Data); err == nil {
			result = append(result, common.Node[T]{
				Data:   value,
				GrpId:  n.GrpID,
				Uuid:   n.Uuid,
				Parent: n.Parent,
				Child:  n.Child,
			})
		}
	}
	return result
}
