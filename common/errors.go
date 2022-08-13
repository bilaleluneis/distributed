/*	@author: Bilal El Uneis
	@since: August 2022
	@email: bilaleluneis@gmail.com	*/

package common

// Error credit to https://dave.cheney.net/2016/04/07/constant-errors
type Error string

func (e Error) Error() string { return string(e) }

const (
	ReqGrpIdErr          = Error("group id is required")
	ReqUuidErr           = Error("uuid is required")
	DoesNotExistGrpIdErr = Error("group id does not exist")
	NoResultsErr         = Error("no results")
	NoWorkerAvailErr     = Error("no available workers")
	MultipleMatchErr     = Error("multiple matches found")
	NonToDelErr          = Error("nothing to delete")
	RpcServiceRegErr     = Error("failed to register rpc service")
	InitWorkerFailed     = Error("failed to init worker")
)
