/*	@author: Bilal El Uneis
	@since: June 2022
	@email: bilaleluneis@gmail.com	*/

package internal

import (
	"bytes"
	"errors"
)

type RpcNode struct{}

var nodes []bytes.Buffer

// Insert the second param here _ *struct{} is equal to passing *void (nothing)
// and is meant to satisfy the RPC service signature requirements.
func (RpcNode) Insert(val []byte, _ *struct{}) error {

	if len(nodes) == 0 {
		nodes = make([]bytes.Buffer, 5, 5)
	}

	data := bytes.Buffer{}
	data.Write(val)
	nodes = append(nodes, data)
	return nil
}

func (RpcNode) Retrieve(remove bool, result *[]byte) error {
	if len(nodes) == 0 {
		return errors.New("nothing to pop")
	}
	*result = nodes[len(nodes)-1].Bytes()
	if remove {
		nodes = nodes[:len(nodes)-1]
	}
	return nil
}
