/*	@author: Bilal El Uneis
	@since: July 2022
	@email: bilaleluneis@gmail.com	*/

// Test gob conversion to and from buffer
package research

import (
	"encoding/gob"
	"testing"
)

func TestSimpleGob(t *testing.T) {

	// convert C{} to bytes via encoder
	enc := gob.NewEncoder(&network)
	err := enc.Encode(C{1, "hello"})
	if err != nil {
		t.Fail()
	}

	// read byts from decoder and convert back to type C
	dec := gob.NewDecoder(&network)
	var c C
	err = dec.Decode(&c)
	if err != nil {
		t.Fail()
	}
	if c.I != 1 && c.getJ() != "hello" {
		t.Fail()
	}

}

func TestGenericGob(t *testing.T) {

	err := toNetwork(Node[int]{1})
	if err != nil {
		t.Fail()
	}

	result, err := fromNetwork[Node[int]]()
	if err != nil {
		t.Fail()
	}
	if result.get() != 1 {
		t.Fail()
	}

}
