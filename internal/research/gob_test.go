/*	@author: Bilal El Uneis
	@since: July 2022
	@email: bilaleluneis@gmail.com	*/

// Test gob conversion to and from buffer
package research

import (
	"bytes"
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

// pay attention to comments bellow to understand
// how to make Interface types work with GOB
func TestInterfaceOverGob(t *testing.T) {
	// must register the concreet type
	gob.Register(TestImpl{})
	var network bytes.Buffer
	var err error

	// assign instance of concreet type with exported Fields to an interface var
	var sentInterface TestInterface = TestImpl{Flag: true}
	// pass interface by ref to encoder
	if err = gob.NewEncoder(&network).Encode(&sentInterface); err == nil {
		var recievedInterface TestInterface
		// pass recieved interface by ref to decoder
		if err = gob.NewDecoder(&network).Decode(&recievedInterface); err == nil {
			if recievedInterface.TestMethod() {
				return // Test success
			}
		}
	}

	// got here, then Test Failed
	t.Fatalf("Failed due to err %s", err.Error())

}
