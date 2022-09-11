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

type List[T any] struct {
	size       int
	identifier common.GRPID
	root       common.UUID
}

func (list List[T]) Identity() common.GRPID {
	return list.identifier
}

func (list List[T]) Len() int {
	return list.size
}

// Push will push new node on top, this allows for
// faster insertion as all is needed is creating
// new node with with child = current root
// and updating current list root to be the
// newly created node
func (list *List[T]) Push(val T) error {
	// empty identifier = list was not correclty created
	if list.identifier == common.EmptyGrpID {
		return common.CollectionNotInitedErr
	}

	var err error
	// this newly created list
	if list.size == 0 {
		node := common.Node[T]{
			Data:  val,
			GrpId: list.identifier,
			Uuid:  list.root,
		}
		if err = UpdateNode[T](node); err == nil {
			list.size++
		}
		return err
	}

	// not newly created list, so push on top
	var uuid common.UUID
	if _, uuid, err = NewNode(val, list.identifier); err == nil {
		node := common.Node[T]{
			Data:  val,
			GrpId: list.identifier,
			Uuid:  uuid,
			Child: list.root,
		}
		if err = UpdateNode[T](node); err == nil {
			list.root = node.Uuid
			list.size++
		}
	}

	return err
}

// NewLinkedList FIXME: do not need this as vargs version takes 0 or more
func NewLinkedList[T any]() (List[T], error) {
	var emptyVal T
	grpId, uuid, err := NewNode(emptyVal, common.EmptyGrpID)
	var linkedList List[T]
	if err == nil {
		linkedList.identifier = grpId
		linkedList.root = uuid
	}
	return linkedList, err
}

func NewLinkedListWithValues[T any](values ...T) (List[T], error) {
	linkedList, err := NewLinkedList[T]()
	for _, value := range values {
		if err = linkedList.Push(value); err != nil {
			return List[T]{}, err
		}
	}
	return linkedList, err
}
