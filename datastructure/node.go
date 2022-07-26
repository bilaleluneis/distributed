/*	@author: Bilal El Uneis
	@since: July 2022
	@email: bilaleluneis@gmail.com	*/

package datastructure

import (
	"bytes"
	"distributed/common"
	"distributed/internal"
	"encoding/gob"
	"errors"
	"net/rpc"
)

type Node[T any] struct {
	Data T
	Prev common.Location
	Next common.Location
}

func Insert[T any](node Node[T], loc common.Location) (common.UUID, error) {
	buffer := bytes.NewBuffer([]byte{})
	err := gob.NewEncoder(buffer).Encode(node)
	var uuid common.UUID
	if err != nil {
		return uuid, err
	}
	var client *rpc.Client
	if client, err = internal.GetTcpClient(loc); err == nil {
		err = client.Call(internal.INSERT, buffer.Bytes(), &uuid)
	}
	return uuid, err
}

func Retrieve[T any](param common.SearchParams) (Node[T], error) {
	var node Node[T]
	client, err := internal.GetTcpClient(param.Address)
	if err != nil {
		return node, errors.New("call to server failed")
	}
	var result []byte
	err = client.Call(internal.RETRIEVE, param, &result)
	if err != nil {
		return node, errors.New("call to RpcNode.Retrieve failed")
	}
	buffer := bytes.NewBuffer(result)
	err = gob.NewDecoder(buffer).Decode(&node)
	return node, err
}

func Update[T any](at common.Location, newValue Node[T]) error {
	node := newValue
	buffer := bytes.NewBuffer([]byte{})
	err := gob.NewEncoder(buffer).Encode(node)
	if err != nil {
		return err
	}
	client, err := internal.GetTcpClient(at)
	if err != nil {
		return err
	}
	updatedNode := internal.UpdatedNode{
		Uuid: at.Uuid,
		Node: buffer.Bytes(),
	}
	err = client.Call(internal.UPDATE, updatedNode, nil)
	return err
}
