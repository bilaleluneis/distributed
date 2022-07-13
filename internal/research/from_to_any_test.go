/*	@author: Bilal El Uneis
	@since: July 2022
	@email: bilaleluneis@gmail.com	*/

// Test conversion from Constraint input to any and cast back
package research

import (
	"testing"
)

func TestSimpleType(t *testing.T) {
	insertData[int](1)
	result, err := retrieveData[int]()
	if err != nil || result != 1 {
		t.Fail()
	}
}

func TestComplexType(t *testing.T) {
	c := []int{1, 2, 3}
	insertData[[]int](c)
	result, err := retrieveData[[]int]()
	if err != nil || result[1] != 2 {
		t.Fail()
	}
}

func TestCustomType(t *testing.T) {
	insertData[C](C{1, "hi"})
	result, err := retrieveData[C]()
	if err != nil || result.getJ() != "hi" {
		t.Fail()
	}
}

func TestCustomGenericType(t *testing.T) {
	insertData[Node[int]](Node[int]{1})
	result, err := retrieveData[Node[int]]()
	if err != nil || result.Value != 1 {
		t.Fail()
	}
}
