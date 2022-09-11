/*	@author: Bilal El Uneis
	@since: August 2022
	@email: bilaleluneis@gmail.com	*/

package datastructure

import (
	"distributed/common"
	"testing"
)

func TestCreation(t *testing.T) {
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

func TestFailedPush(t *testing.T) {
	l := List[int]{}
	if err := l.Push(1); err != common.CollectionNotInitedErr {
		t.Fail()
	}
}

// NOTE: before debugging failure in this test
// make sure TestLinkedListCreation passes
func TestSuccessfulPush(t *testing.T) {
	l, _ := NewLinkedList[int]()
	if err := l.Push(1); err != nil {
		t.Fatalf("push with error %s", err.Error())
	}
	if l.Len() != 1 {
		t.Fatalf("size after push should be 1 instead got %d", l.Len())
	}
}
