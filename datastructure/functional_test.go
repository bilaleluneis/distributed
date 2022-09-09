/*	@author: Bilal El Uneis
	@since: Sep 2022
	@email: bilaleluneis@gmail.com	*/

package datastructure

import (
	"distributed"
	"distributed/common"
	"distributed/internal"
	"encoding/gob"
	"testing"
)

// FIXME: remove init once functional types are handeled diff for GOB
func init() {
	gob.Register(internal.Compute{})
	distributed.RegisterFilter[int](valueFilter{})
	distributed.RegisterMapper[int, int](valueMapper{})
	distributed.RegisterReducer[int, int](sumReducer{})
}

func TestCompute(t *testing.T) {
	values := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	l, _ := NewLinkedListWithValues(values...)
	if result, err := Compute[int](l); err == nil {
		if len(result) != len(values) {
			t.Fatalf("size check")
		}
		for _, n := range result {
			if count, exist := common.Contains[int](values, n); !exist || count != 1 {
				t.Fatalf("exist and/or count check")
			}
		}
	} else {
		t.Fatalf("compute returned error %s", err.Error())
	}
}

type sumReducer struct{}

func (sumReducer) Reduce(n ...common.NodeLike[int]) int {
	sum := 0
	for _, item := range n {
		sum += item.GetData()
	}
	return sum
}

func TestReduce(t *testing.T) {
	values := []int{1, 2, 3, 4, 5}
	l, _ := NewLinkedListWithValues(values...)
	finalReduction := func(in []int) int {
		sum := 0
		for _, i := range in {
			sum += i
		}
		return sum
	}
	result, err := Reduce[int, int](l, sumReducer{}, finalReduction)
	if err != nil {
		t.Fatalf("reduction failed %s", err.Error())
	}
	if result != 15 {
		t.Fatalf("wrong reduction value")
	}
}

// NOTE: Types Filter/Map must have exported Fileds or will not work with GOB
type valueFilter struct{ Value int }

func (v valueFilter) Filter(n common.NodeLike[int]) bool {
	return n.GetData() == v.Value
}

func TestFilterCompute(t *testing.T) {
	values := []int{1, 2, 3, 4, 5, 2, 7, 8, 9, 2}
	l, _ := NewLinkedListWithValues(values...)
	if err := Filter[int](l, valueFilter{2}); err == nil {
		result, err := Compute[int](l)
		if err != nil {
			t.Fatalf("compute falied %s", err.Error())
		}
		if len(result) != 3 {
			t.Fatalf("wrong result count")
		}
		if count, exist := common.Contains(result, 2); !exist || count != 3 {
			t.Fatalf("got wrong occurance of value 2")
		}
	} else {
		t.Fatalf("filter %s", err.Error())
	}
}

type valueMapper struct{ Value int }

func (m valueMapper) Map(n common.NodeLike[int]) common.Node[int] {
	return common.Node[int]{
		Data:   m.Value,
		GrpId:  n.GetGrpID(),
		Uuid:   n.GetUuID(),
		Parent: n.GetParent(),
		Child:  n.GetChild(),
	}
}

func TestMapCompute(t *testing.T) {
	values := []int{1, 2, 3, 4, 5, 2, 7, 8, 9, 2}
	l, _ := NewLinkedListWithValues(values...)
	if err := Map[int, int](l, valueMapper{2}); err == nil {
		result, err := Compute[int](l)
		if err != nil {
			t.Fatalf("compute %s", err.Error())
		}
		if len(result) != len(values) {
			t.Fatalf("wrong result length")
		}
		if count, exist := common.Contains(result, 2); !exist || count != len(values) {
			t.Fatalf("got wrong occurance of mapped value 2")
		}
	} else {
		t.Fatalf("map %s", err.Error())
	}
}
