/*	@author: Bilal El Uneis
	@since: Nov 2022
	@email: bilaleluneis@gmail.com	*/

package research

import (
	"bytes"
	"distributed/common"
	"encoding/gob"
	"testing"
)

func TestSelfEncodedDecodedType(t *testing.T) {
	startTestWorker(SimpleRpc{}, t)
	regedWorker := common.GetAvailRegWorkers()[0] // get first worker
	reply := SelfEncodedDecoded{}
	if err := regedWorker.Invoke("SimpleRpc.Mutate", SelfEncodedDecoded{i: 1, j: true}, &reply); err != nil {
		t.Fatalf("invoking mutate rpc call failed due %s", err.Error())
	}
	if reply.i != 1 && !reply.j {
		t.Fail()
	}
}

type SimpleRpc struct{}

func (r SimpleRpc) ServiceName() common.Service {
	return "SimpleRpc"
}

func (SimpleRpc) Mutate(from SelfEncodedDecoded, to *SelfEncodedDecoded) error {
	*to = SelfEncodedDecoded{
		i: from.i,
		j: from.j,
	}
	return nil
}

// SelfEncodedDecoded : will implement GobIncoder and GobDecoder
// which will allow to break free from GOB rules, like un-exported fields
type SelfEncodedDecoded struct {
	i int
	j bool
}

func (sed SelfEncodedDecoded) GobEncode() ([]byte, error) {
	var network bytes.Buffer
	encoder := gob.NewEncoder(&network)
	_ = encoder.Encode(sed.i)
	_ = encoder.Encode(sed.j)
	return network.Bytes(), nil
}

func (sed *SelfEncodedDecoded) GobDecode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	*sed = SelfEncodedDecoded{}
	_ = decoder.Decode(&sed.i)
	_ = decoder.Decode(&sed.j)
	return nil
}
