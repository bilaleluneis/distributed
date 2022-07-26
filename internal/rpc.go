/*	@author: Bilal El Uneis
	@since: June 2022
	@email: bilaleluneis@gmail.com	*/

package internal

import (
	"bytes"
	"distributed/common"
	"errors"
)

// available RPC calls
const (
	INSERT   string = "RpcNode.Insert"
	RETRIEVE        = "RpcNode.Retrieve"
	UPDATE          = "RpcNode.Update"
)

// RpcNode : RPC service type
type RpcNode struct{}

var nodes = make(map[common.UUID]*bytes.Buffer)

func (RpcNode) Insert(node []byte, uuid *common.UUID) error {
	_data := bytes.NewBuffer(node)
	*uuid = common.GenUUID()
	for _, ok := nodes[*uuid]; ok; {
		*uuid = common.GenUUID()
	}
	nodes[*uuid] = _data
	return nil
}

func (RpcNode) Retrieve(params common.SearchParams, result *[]byte) error {
	if len(nodes) == 0 || params.Address.Uuid == "" {
		return errors.New("nothing to pop")
	}
	if _, ok := nodes[params.Address.Uuid]; !ok {
		return errors.New("uuid does not exist")
	}
	*result = nodes[params.Address.Uuid].Bytes()
	if params.Remove {
		delete(nodes, params.Address.Uuid)
	}
	return nil
}

// UpdatedNode Update the second empty struct param is = to passing void
type UpdatedNode struct {
	Uuid common.UUID
	Node []byte
}

func (RpcNode) Update(data UpdatedNode, _ *struct{}) error {
	updatedNode := bytes.NewBuffer(data.Node)
	if _, ok := nodes[data.Uuid]; !ok {
		return errors.New("uuid does not exist")
	}
	nodes[data.Uuid] = updatedNode
	return nil
}
