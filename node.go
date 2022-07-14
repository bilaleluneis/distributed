/*	@author: Bilal El Uneis
	@since: July 2022
	@email: bilaleluneis@gmail.com	*/

package distributed

import (
	"bytes"
	"encoding/gob"
	"errors"
)

type Node[T any] struct{}

func (n Node[T]) Push(val T) error {
	buffer := bytes.NewBuffer([]byte{})
	err := gob.NewEncoder(buffer).Encode(val)
	if err != nil {
		return err
	}
	err = client.Call("RpcNode.Insert", buffer.Bytes(), nil)
	return err
}

func (n Node[T]) Pop(remove bool) (T, error) {
	var data []byte
	err := client.Call("RpcNode.Retrieve", remove, &data)

	var result T
	if err != nil {
		return result, errors.New("call to server failed")
	}

	buffer := bytes.NewBuffer([]byte{})
	buffer.Write(data)
	dec := gob.NewDecoder(buffer)
	err = dec.Decode(&result)
	if err != nil {
		return result, errors.New("node pop() conversion failure")
	}

	return result, nil
}
