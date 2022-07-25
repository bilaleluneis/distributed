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
	//TODO: check if uuid exist then gen another one
	*uuid = common.GenUUID()
	nodes[*uuid] = _data
	return nil
}

// Retrieve TODO: refactor to account for when Uuid is not provided or not present
// TODO: define error constants in common package
func (RpcNode) Retrieve(params common.SearchParams, result *[]byte) error {
	if len(nodes) == 0 {
		return errors.New("nothing to pop")
	}
	//TODO: if Uuid is provided check if exist in nodes else return error
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
	nodes[data.Uuid] = updatedNode
	return nil
}
