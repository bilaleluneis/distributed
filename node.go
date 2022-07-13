/*	@author: Bilal El Uneis
	@since: July 2022
	@email: bilaleluneis@gmail.com	*/

package distributed

import (
	"bytes"
	"encoding/gob"
	"errors"
)

type Node struct{}

func (n Node) Push(val int) error {
	buffer := bytes.NewBuffer([]byte{})
	err := gob.NewEncoder(buffer).Encode(val)
	if err != nil {
		return err
	}
	err = client.Call("RpcNode.Insert", buffer.Bytes(), nil)
	return err
}

func (n Node) Pop(remove bool) (int, error) {
	var data []byte
	err := client.Call("RpcNode.Retrieve", remove, &data)

	var result int
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
