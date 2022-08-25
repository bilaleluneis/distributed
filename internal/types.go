/*	@author: Bilal El Uneis
	@since: August 2022
	@email: bilaleluneis@gmail.com	*/

package internal

import "distributed/common"

type Functional interface {
	Func([]RpcNode) []RpcNode
}

type FunctionParam struct {
	GrpID    common.GRPID
	Function Functional
}
