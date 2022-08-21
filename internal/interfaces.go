/*	@author: Bilal El Uneis
	@since: August 2022
	@email: bilaleluneis@gmail.com	*/

package internal

import "distributed/common"

type Filterer interface {
	Filter(map[common.GRPID][]RpcNode) []RpcNode
}
