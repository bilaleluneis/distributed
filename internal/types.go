/*	@author: Bilal El Uneis
	@since: June 2022
	@email: bilaleluneis@gmail.com	*/

package internal

import (
	"errors"
)

type RpcNode struct{}

var nodes []int

func (RpcNode) Insert(val int, _ *int) error {

	if len(nodes) == 0 {
		nodes = make([]int, 5, 5)
	}
	nodes = append(nodes, val)
	return nil
}

func (RpcNode) Retrieve(remove bool, result *int) error {
	if len(nodes) == 0 {
		return errors.New("nothing to pop")
	}
	*result = nodes[len(nodes)-1]
	if remove {
		nodes = nodes[:len(nodes)-1]
	}
	return nil
}
