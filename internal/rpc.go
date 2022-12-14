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
	DELAYED    = common.Service("RpcNodeService.Delayed")
	IMMEDIATE  = common.Service("RpcNodeService.Immediate")
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
// TODO: mutex and sync access ?
type RpcNodeService struct {
	Nodes map[common.GRPID][]RpcNode
	Ops   map[common.GRPID][]FunctionalOp
}

// ServiceName implement common.ServiceProvider
func (RpcNodeService) ServiceName() common.Service {
	return "RpcNodeService"
}

// New returns an available Group ID on this worker
// the client is responsible for making sure that this
// Group I D does not already exist on other workers
func (rns *RpcNodeService) New(_ common.NONE, grpId *common.GRPID) error {
	id := common.GenUUID()
	for _, ok := rns.Nodes[id]; ok; {
		id = common.GenUUID()
	}
	*grpId = id
	common.Log.Debug("new group with value %s created", id)
	return nil
}

func (rns RpcNodeService) GrpIdExist(grpId common.GRPID, exist *bool) error {
	if _, ok := rns.Nodes[grpId]; ok {
		*exist = true
	}
	return nil
}

func (rns RpcNodeService) UuidExist(node RpcNode, exist *bool) error {
	if node.GrpID == common.EmptyGrpID {
		return common.ReqGrpIdErr
	}
	if node.Uuid == common.EmptyUUID {
		return common.ReqUuidErr
	}
	grpId := node.GrpID
	uuid := node.Uuid
	if nodes, ok := rns.Nodes[grpId]; ok {
		for _, node := range nodes {
			if node.Uuid == uuid {
				*exist = true
				break
			}
		}
		return nil
	}
	return common.DoesNotExistGrpIdErr
}

// Insert client is responsible for passing both Group ID and UUID
/*
this approach can be used to do both insert , update and duplicate
client will need to check for existance of group id and uuid on
workers and determine if insert or update on worker if uuid exist
or even duplicate on another worker. if both group id and uuid do not exist
this method will generate both.
*/
func (rns *RpcNodeService) Insert(node RpcNode, _ *common.NONE) error {
	if node.GrpID == common.EmptyGrpID {
		return common.ReqGrpIdErr
	}
	if node.Uuid == common.EmptyUUID {
		return common.ReqUuidErr
	}
	grpId, uuid := node.GrpID, node.Uuid
	for index, currNode := range rns.Nodes[grpId] {
		if currNode.GrpID == grpId && currNode.Uuid == uuid {
			rns.Nodes[grpId][index] = node
			return nil
		}
	}
	rns.Nodes[grpId] = append(rns.Nodes[grpId], node)
	common.Log.Debug("inserted Node for group %s with uuid %s", grpId, uuid)
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
	if grpId == common.EmptyGrpID {
		return common.ReqGrpIdErr
	}
	updatedNodeList := make([]RpcNode, 0)
	for _, node := range rns.Nodes[grpId] {
		if !found(node, nodesToDel) {
			updatedNodeList = append(updatedNodeList, node)
		}

	}
	rns.Nodes[grpId] = updatedNodeList
	return nil
}

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
func (rns RpcNodeService) Retrieve(criteria RpcNode, result *[]RpcNode) error {
	if criteria.GrpID == common.EmptyGrpID {
		return common.ReqGrpIdErr
	}

	rnodes, ok := rns.Nodes[criteria.GrpID]
	if !ok {
		return common.DoesNotExistGrpIdErr
	}
	if len(rnodes) == 0 {
		return common.NoResultsErr
	}

	resultsFound := make([]RpcNode, 0)
	for _, n := range rnodes {
		found := false
		if criteria.Uuid != common.EmptyUUID && n.Uuid == criteria.Uuid {
			found = true
		}
		if len(criteria.Data) != 0 && bytes.Equal(criteria.Data, n.Data) {
			found = true
		}
		if found {
			resultsFound = append(resultsFound, n)
		}
	}

	if len(resultsFound) == 0 {
		return common.NoResultsErr
	}

	*result = resultsFound
	return nil
}

// Delayed will store functional operation but not evaluate (lazy eval)
func (rns *RpcNodeService) Delayed(by FuncParam, _ *common.NONE) error {
	if by.GrpId == common.EmptyGrpID {
		return common.ReqGrpIdErr
	}
	grpId := by.GrpId
	_, ok := rns.Ops[grpId]
	if ok {
		rns.Ops[grpId] = append(rns.Ops[grpId], by.Op)
	} else {
		rns.Ops[grpId] = []FunctionalOp{by.Op}
	}
	common.Log.Debug("inserted delay op for group %s", grpId)
	return nil
}

// Immediate will eager eval functional operations store if any then
// evaluate functional operation passed and return result
func (rns *RpcNodeService) Immediate(by FuncParam, result *[]RpcNode) error {
	grpId := by.GrpId
	var validGrpId bool
	var err error
	if _ = rns.GrpIdExist(grpId, &validGrpId); validGrpId {
		currData := rns.Nodes[grpId]
		// evaluate previously delayed Ops if any
		if ops, ok := rns.Ops[grpId]; ok {
			defer delete(rns.Ops, grpId) // clear all stored Ops once done
			for _, f := range ops {
				currData = f.Eval(currData)
			}
		}
		// evaluate current operation and store results as bytes
		*result = by.Op.Eval(currData)
		common.Log.Debug("executing immediate op for group %s", grpId)
		return err
	}
	return common.DoesNotExistGrpIdErr
}
