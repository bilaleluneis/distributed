/*	@author: Bilal El Uneis
	@since: June 2022
	@email: bilaleluneis@gmail.com	*/

package internal

import (
	"bytes"
	"distributed/common"
)

// available RPC calls
const (
	NEW        = common.Service("RpcNodeService.New")
	INSERT     = common.Service("RpcNodeService.Insert")
	RETRIEVE   = common.Service("RpcNodeService.Retrieve")
	DELETE     = common.Service("RpcNodeService.Delete")
	FILTER     = common.Service("RpcNodeService.Filter")
	GRPIDEXIST = common.Service("RpcNodeService.GrpIdExist")
)

// RpcNode will hold information about node and will be used as DTO back and forth
type RpcNode struct {
	Data   []byte
	GrpID  common.GRPID // which group this node belongs to
	Uuid   common.UUID
	Parent common.UUID // previous node
	Child  common.UUID // next node
}

// RpcNodeService : RPC service type
type RpcNodeService struct {
	nodes map[common.GRPID][]RpcNode
}

// New returns an available Group ID on this worker
// the client is responsible for making sure that this
// Group ID does not already exist on other workers
// TODO: might not need this method if I do GRPID generation on client side
func (rns *RpcNodeService) New(_ common.NONE, grpId *common.GRPID) error {
	id := common.GenUUID()
	for _, ok := rns.nodes[id]; ok; {
		id = common.GenUUID()
	}
	*grpId = id
	return nil
}

func (rns RpcNodeService) GrpIdExist(grpId common.GRPID, exist *bool) error {
	if _, ok := rns.nodes[grpId]; ok {
		*exist = true
	}
	return nil
}

// TODO: create UuidExist(Node, exist *bool) err
// use node to pass group id and generated uuid to check

// Insert client is responsible for passing both Group ID and UUID
/*
this approach can be used to do both insert , update and duplicate
client will need to check for existance of group id and uuid on
workers and determine if insert or update on worker if uuid exist
or even duplicate on another worker. if both group id and uuid do not exist
this method will generate both.
*/
func (rns *RpcNodeService) Insert(node RpcNode, _ *common.NONE) error {
	if node.GrpID == "" {
		return common.ReqGrpIdErr
	}
	if node.Uuid == "" {
		return common.ReqUuidErr
	}
	grpId, uuid := node.GrpID, node.Uuid
	//TODO: can I use channels and goroutine here?
	for index, currNode := range rns.nodes[grpId] {
		if currNode.GrpID == grpId && currNode.Uuid == uuid {
			rns.nodes[grpId][index] = node
			return nil
		}
	}
	rns.nodes[grpId] = append(rns.nodes[grpId], node)
	return nil
}

// Delete will use criteria match and delete one or more that match
// criteria.GrpID is mandatory and should be same for all
// will grab GrpID from first element and reuse while looping
// client is responsible for passing correct info slice with uniform GrpID
// and valid / existing UUID withen that Group
// invalid and non found UUID withen Grp will be ignored silently
func (rns *RpcNodeService) Delete(nodesToDel []RpcNode, _ *common.NONE) error {
	if len(nodesToDel) == 0 {
		return common.NonToDelErr
	}
	grpId := nodesToDel[0].GrpID
	if grpId == "" {
		return common.ReqGrpIdErr
	}
	updatedNodeList := make([]RpcNode, 0)
	for _, node := range rns.nodes[grpId] {
		if !found(node, nodesToDel) {
			updatedNodeList = append(updatedNodeList, node)
		}

	}
	rns.nodes[grpId] = updatedNodeList
	return nil
}

// TODO: might want to use go routines if len(from) is larger than some value
func found(node RpcNode, from []RpcNode) bool {
	for _, n := range from {
		if node.GrpID == n.GrpID && node.Uuid == n.Uuid {
			return true
		}
	}
	return false
}

// Retrieve will use criteria parameter to match as many as possible
// and return slice with results matching criteria
// criteria.GrpID is manadatory, the rest are optional
// TODO: consider using go routines to search faster
// FIXME: for now just checking against data and uuid, impl is slow
func (rns RpcNodeService) Retrieve(criteria RpcNode, result *[]RpcNode) error {
	if criteria.GrpID == "" {
		return common.ReqGrpIdErr
	}

	rnodes, ok := rns.nodes[criteria.GrpID]
	if !ok {
		return common.DoesNotExistGrpIdErr
	}
	if len(rnodes) == 0 {
		return common.NoResultsErr
	}

	resultsFound := make([]RpcNode, 0)
	for _, n := range rnodes {
		found := false
		if criteria.Uuid != "" && n.Uuid == criteria.Uuid {
			found = true
		}
		if len(criteria.Data) != 0 && bytes.Equal(criteria.Data, n.Data) {
			found = true
		}
		if found {
			resultsFound = append(resultsFound, n)
		}
	}

	*result = resultsFound
	return nil
}

func (rns RpcNodeService) Filter(by Filterer, result *[]RpcNode) error {
	*result = by.Filter(rns.nodes)
	return nil
}
