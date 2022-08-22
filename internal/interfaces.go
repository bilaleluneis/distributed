/*	@author: Bilal El Uneis
	@since: August 2022
	@email: bilaleluneis@gmail.com	*/

package internal

import "distributed/common"

type Grouped interface {
	ForGroup() common.GRPID
}

type UniqueIdentifiable interface {
	WithUUID() common.UUID
}

type Filterer interface {
	Grouped
	Filter([]RpcNode) []RpcNode
}
