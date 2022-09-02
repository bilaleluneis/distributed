/*	@author: Bilal El Uneis
	@since: August 2022
	@email: bilaleluneis@gmail.com	*/

package internal

import (
	"bytes"
	"distributed/common"
	"encoding/gob"
)

type FunctionalOp interface {
	Eval([]RpcNode) []RpcNode
}

type FuncParam struct {
	Op    FunctionalOp
	GrpId common.GRPID
}

type Filter[F common.Filterer[T], T any] struct {
	WithFilter F
}

func (f Filter[F, T]) Eval(rpcNodes []RpcNode) []RpcNode {
	result := make([]RpcNode, 0)
	for _, rpcNode := range rpcNodes {
		n := decode[T](rpcNode)
		var node common.NodeLike[T] = &n[0]
		if f.WithFilter.Filter(node) {
			result = append(result, rpcNode)
		}
	}
	return result
}

type Reduce[R common.Reducer[T, O], T any, O any] struct {
	WithReducer R
}

func (r Reduce[R, T, O]) Eval(rpcNodes []RpcNode) []RpcNode {
	nodes := decode[T](rpcNodes...)
	reductionList := make([]common.NodeLike[T], 0)
	for _, n := range nodes {
		var reductionNode common.NodeLike[T] = &n
		reductionList = append(reductionList, reductionNode)
	}
	return []RpcNode{encode[O](&common.Node[O]{
		Data: r.WithReducer.Reduce(reductionList...),
	})[0]}
}

type Map[M common.Mapper[T], T any] struct {
	WithMapper M
}

func (m Map[M, T]) Eval(rpcNodes []RpcNode) []RpcNode {
	result := make([]RpcNode, 0)
	for _, rpcNode := range rpcNodes {
		node := decode[T](rpcNode)[0]
		mappedNode := m.WithMapper.Map(&node)
		result = append(result, encode[T](mappedNode)[0])
	}
	return result
}

func encode[T any](nodes ...common.NodeLike[T]) []RpcNode {
	result := make([]RpcNode, 0)
	for _, n := range nodes {
		var buffer bytes.Buffer
		err := gob.NewEncoder(&buffer).Encode(n.GetData())
		if err == nil {
			result = append(result, RpcNode{
				Data:   buffer.Bytes(),
				GrpID:  n.GetGrpID(),
				Uuid:   n.GetUuID(),
				Parent: n.GetParent(),
				Child:  n.GetChild(),
			})
		}
	}
	return result
}

func decode[T any](rpcNodes ...RpcNode) []common.Node[T] {
	result := make([]common.Node[T], 0)
	for _, n := range rpcNodes {
		var buffer bytes.Buffer
		buffer.Write(n.Data)
		var value T
		err := gob.NewDecoder(&buffer).Decode(&value)
		if err == nil {
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
