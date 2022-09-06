/*	@author: Bilal El Uneis
	@since: August 2022
	@email: bilaleluneis@gmail.com	*/

package datastructure

import (
	"distributed/common"
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
}

func (ll LinkedList[T]) Identity() common.GRPID {
	return ll.identifier
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
		node := common.Node[T]{
			Data:  val,
			GrpId: ll.identifier,
			Uuid:  ll.root,
		}
		if err = UpdateNode[T](node); err == nil {
			ll.size++
		}
		return err
	}

	// not newly created list, so push on top
	var uuid common.UUID
	if _, uuid, err = NewNode(val, ll.identifier); err == nil {
		node := common.Node[T]{
			Data:  val,
			GrpId: ll.identifier,
			Uuid:  uuid,
			Child: ll.root,
		}
		if err = UpdateNode[T](node); err == nil {
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
