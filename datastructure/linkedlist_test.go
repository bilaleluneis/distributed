/*	@author: Bilal El Uneis
	@since: August 2022
	@email: bilaleluneis@gmail.com	*/

package datastructure

import (
	"distributed/common"
	"distributed/internal"
	"encoding/gob"
	"testing"
)

func init() {
	gob.Register(lessThanHundredFilter{})
	gob.Register(accumilateReducer{})
}

func TestLinkedListCreation(t *testing.T) {
	// creating an Empty list
	list, err := NewLinkedList[int]()
	if err != nil {
		t.Fatal("creating new linked list failed")
	}
	if list.Len() != 0 {
		t.Fatalf("expected list size to be 0 got %d", list.Len())
	}

	// creating list with initial values
	list, err = NewLinkedListWithValues(1, 2, 3, 4)
	if err != nil {
		t.Fatal("failed creating list with initial values")
	}
	if list.Len() != 4 {
		t.Fatalf("expected list size to be 3 got %d", list.Len())
	}
}

func TestLinkedListFailedPush(t *testing.T) {
	l := LinkedList[int]{}
	if err := l.Push(1); err != common.CollectionNotInitedErr {
		t.Fail()
	}
}

func TestLinkedListPush(t *testing.T) {
	var err error
	var l LinkedList[int]

	// create new linked list
	l, err = NewLinkedList[int]()
	if err != nil {
		t.Fatalf("creating new linked list failure")
	}
	if l.Len() != 0 {
		t.Fatalf("new linked list size is not zero, is %d", l.Len())
	}

	// first push
	if err = l.Push(1); err != nil {
		t.Fatalf("first push failed")
	}
	if l.Len() != 1 {
		t.Fatalf("size after first push should be 1 instead %d", l.Len())
	}

	// second push
	if err = l.Push(2); err != nil {
		t.Fatalf("second push failed")
	}
	if l.Len() != 2 {
		t.Fatalf("size after second push should be 2 instead %d", l.Len())
	}
}

type lessThanHundredFilter struct{}

func (lessThanHundredFilter) Func(nodes []internal.RpcNode) []internal.RpcNode {
	result := make([]internal.RpcNode, 0)
	if len(nodes) != 0 {
		for index, n := range nodes {
			if decode[int]([]internal.RpcNode{n})[0].Data < 100 {
				result = append(result, nodes[index])
			}
		}
	}
	return result
}

func TestLinkedListFilter(t *testing.T) {
	if l, err := NewLinkedListWithValues[int](1, 2, 3, 400, 99, 500); err == nil {
		if result, err := l.Filter(lessThanHundredFilter{}).Compute(); err == nil {
			if len(result) == 4 {
				return
			}
		}
	}
	t.Fail()
}

type accumilateReducer struct{}

func (accumilateReducer) Func(nodes []internal.RpcNode) []internal.RpcNode {
	var result int
	if len(nodes) != 0 {
		for _, n := range nodes {
			result = result + decode[int]([]internal.RpcNode{n})[0].Data
		}
	}
	return encode[int]([]Node[int]{{Data: result}})
}

func TestLinkedListReduce(t *testing.T) {
	values := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 100, 1000}
	if l, err := NewLinkedListWithValues[int](values...); err == nil {
		if result, err := l.Filter(lessThanHundredFilter{}).Reduce(accumilateReducer{}); err == nil {
			if result == 55 {
				return
			}
		}
		t.Fail()
	}
}
