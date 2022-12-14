/*	@author: Bilal El Uneis
	@since: June 2022
	@email: bilaleluneis@gmail.com	*/

package common

import (
	"fmt"
	"math/rand"
	"net/rpc"
	"time"
)

type UUID = string
type GRPID = string
type Service = string
type NONE = struct{}

const (
	EmptyGrpID = GRPID("")
	EmptyUUID  = UUID("")
)

// ServiceProvider will be used as constraint
// for Service like types
type ServiceProvider interface {
	ServiceName() Service
}

type NodeLike[T any] interface {
	GetData() T
	GetGrpID() GRPID
	GetUuID() UUID
	GetParent() UUID
	GetChild() UUID
}

type Collection interface {
	Len() int
	Identity() GRPID
}

type Node[T any] struct {
	Data   T
	GrpId  GRPID
	Uuid   UUID
	Parent UUID
	Child  UUID
}

func (n Node[T]) GetData() T      { return n.Data }
func (n Node[T]) GetGrpID() GRPID { return n.GrpId }
func (n Node[T]) GetUuID() UUID   { return n.Uuid }
func (n Node[T]) GetParent() UUID { return n.Parent }
func (n Node[T]) GetChild() UUID  { return n.Child }
func (n Node[T]) String() string {
	rep := "Node[T]: grpId=%s uuid=%s parent=%s child=%s"
	return fmt.Sprintf(rep, n.GrpId, n.Uuid, n.Parent, n.Child)
}

type Filterer[T any] interface {
	Filter(NodeLike[T]) bool
}

type Mapper[T any, R any] interface {
	Map(NodeLike[T]) Node[R]
}

type Reducer[T any, R any] interface {
	Reduce(...NodeLike[T]) R
}

// RegisteredWorker
// TODO: implement asyncInvoke
type RegisteredWorker struct {
	host   string
	port   int
	inited bool //indicator of initalization via regestration
}

func (w RegisteredWorker) Invoke(s Service, args any, result any) error {
	var err error = WorkerNotValidErr
	if w.inited {
		address := fmt.Sprintf("%s:%d", w.host, w.port)
		var client *rpc.Client
		client, err = rpc.Dial("tcp", address)
		if err != nil {
			return err
		}
		defer func(client *rpc.Client) {
			err = client.Close()
		}(client)
		err = client.Call(s, args, result)
	}
	return err
}

var registeredWorkers = make([]RegisteredWorker, 0)

// RegisterWorker performs initalization of registeredWorkers once
func RegisterWorker(host string, port int) {
	registeredWorkers = append(registeredWorkers, RegisteredWorker{
		host:   host,
		port:   port,
		inited: true,
	})
}

func GetAvailRegWorkers() []RegisteredWorker {
	if len(registeredWorkers) == 0 {
		return []RegisteredWorker{}
	}
	copyOfWorkers := make([]RegisteredWorker, 0)
	for _, w := range registeredWorkers {
		copyOfWorkers = append(copyOfWorkers, w)
	}
	return copyOfWorkers
}

func GetRandomAvailRegWorker() (RegisteredWorker, error) {
	if len(registeredWorkers) == 0 {
		return RegisteredWorker{}, NoWorkerAvailErr
	}
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	index := r.Intn(len(registeredWorkers))
	return GetAvailRegWorkers()[index], nil
}
