/*	@author: Bilal El Uneis
	@since: August 2022
	@email: bilaleluneis@gmail.com	*/

package datastructure

import (
	"distributed/common"
	"distributed/internal"
	"encoding/gob"
)

const (
	FILTEROPERR  = common.Error("filter operation failure")
	MAPOPERR     = common.Error("map operation failure")
	REDUCEOPERR  = common.Error("reduce operation failure")
	COMPUTEOPERR = common.Error("compute operation failure")
)

type LinkedList[T any] struct {
	size       int
	identifier common.GRPID
	root       common.UUID
	err        error
}

func (ll LinkedList[T]) Len() int {
	return ll.size
}

// Push will push new node on top, this allows for
// faster insertion as all is needed is creating
// new node with with child = current root
// and updating current list root to be the
// newly created node
func (ll *LinkedList[T]) Push(val T) error {
	// empty identifier = list was not correclty created
	if ll.identifier == common.EmptyGrpID {
		return common.CollectionNotInitedErr
	}

	var err error
	// this newly created list
	if ll.size == 0 {
		node := Node[T]{
			Data:  val,
			GrpId: ll.identifier,
			Uuid:  ll.root,
		}
		if err = UpdateNode(node); err == nil {
			ll.size++
		}
		return err
	}

	// not newly created list, so push on top
	var uuid common.UUID
	if _, uuid, err = NewNode(val, ll.identifier); err == nil {
		node := Node[T]{
			Data:  val,
			GrpId: ll.identifier,
			Uuid:  uuid,
			Child: ll.root,
		}
		if err = UpdateNode(node); err == nil {
			ll.root = node.Uuid
			ll.size++
		}
	}

	return err
}

func NewLinkedList[T any]() (LinkedList[T], error) {
	var emptyVal T
	grpId, uuid, err := NewNode(emptyVal, common.EmptyGrpID)
	var linkedList LinkedList[T]
	if err == nil {
		linkedList.identifier = grpId
		linkedList.root = uuid
	}
	return linkedList, err
}

func NewLinkedListWithValues[T any](values ...T) (LinkedList[T], error) {
	linkedList, err := NewLinkedList[T]()
	for _, value := range values {
		if err = linkedList.Push(value); err != nil {
			return LinkedList[T]{}, err
		}
	}
	return linkedList, err
}

// TODO: need to make it easier to register Functional impls
// TODO: stream line encoding and decoding as part of Functional

// Filter TODO: proper implement
func (ll *LinkedList[T]) Filter(f internal.Functional) *LinkedList[T] {
	if ll.err != nil {
		return ll
	}
	param := internal.FunctionParam{
		GrpID:    ll.identifier,
		Function: f,
	}
	workers := common.GetAvailRegWorkers()
	for _, worker := range workers {
		if err := worker.Invoke(internal.FILTER, &param, &common.NONE{}); err != nil {
			ll.err = FILTEROPERR
			break
		}
	}
	return ll
}

type compute struct{}

func (compute) Func(nodes []internal.RpcNode) []internal.RpcNode {
	return nodes
}

func (ll LinkedList[T]) Compute() ([]Node[T], error) {
	if ll.err != nil {
		return []Node[T]{}, ll.err
	}

	param := internal.FunctionParam{
		GrpID:    ll.identifier,
		Function: &compute{},
	}
	gob.Register(compute{})
	workers := common.GetAvailRegWorkers()
	computeResult := make([]internal.RpcNode, 0)

	for _, worker := range workers {
		var tmpCompute []internal.RpcNode
		if err := worker.Invoke(internal.REDUCE, &param, &tmpCompute); err != nil {
			return []Node[T]{}, err
		}
		if len(tmpCompute) != 0 {
			computeResult = append(computeResult, tmpCompute...)
		}
	}

	if len(computeResult) == 0 {
		return []Node[T]{}, COMPUTEOPERR
	}
	return decode[T](computeResult), nil
}

// Reduce TODO: maybe need to return error or indicator that we got nothing
func (ll LinkedList[T]) Reduce(f internal.Functional) (T, error) {
	var result T
	if ll.err != nil {
		return result, ll.err
	}

	workers := common.GetAvailRegWorkers()
	reduceResult := make([]internal.RpcNode, 0)
	param := internal.FunctionParam{
		GrpID:    ll.identifier,
		Function: f,
	}

	for _, worker := range workers {
		var tmpReduce []internal.RpcNode
		if err := worker.Invoke(internal.REDUCE, &param, &tmpReduce); err != nil {
			return result, REDUCEOPERR
		}
		if len(tmpReduce) != 0 {
			reduceResult = append(reduceResult, tmpReduce[0])
		}
	}

	// after results gathered do one more reduce on client side
	if len(reduceResult) != 0 {
		reduceResult = param.Function.Func(reduceResult)
		result = decode[T](reduceResult)[0].Data
		return result, nil
	}
	return result, REDUCEOPERR
}
