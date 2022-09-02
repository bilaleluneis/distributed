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

type NodeLike[T any] interface {
	GetData() T
	GetGrpID() GRPID
	GetUuID() UUID
	GetParent() UUID
	GetChild() UUID
}

type Node[T any] struct {
	Data   T
	GrpId  GRPID
	Uuid   UUID
	Parent UUID
	Child  UUID
}

func (in Node[T]) GetData() T      { return in.Data }
func (in Node[T]) GetGrpID() GRPID { return in.GrpId }
func (in Node[T]) GetUuID() UUID   { return in.Uuid }
func (in Node[T]) GetParent() UUID { return in.Parent }
func (in Node[T]) GetChild() UUID  { return in.Child }

type Filterer[T any] interface {
	Filter(NodeLike[T]) bool
}

type Mapper[T any] interface {
	Map(NodeLike[T]) Node[T]
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
