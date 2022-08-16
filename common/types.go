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

// TODO refactor API to use this instead of ""
const (
	EmptyGrpID = GRPID("")
	EmptyUUID  = UUID("")
)

// RegisteredWorker
// TODO: implement asyncInvoke
type RegisteredWorker struct {
	host   string
	port   int
	inited bool //indicator of initalization via regestration
}

func (w RegisteredWorker) Invoke(s Service, args any, result any) error {
	//FIXME: error if inited is false
	address := fmt.Sprintf("%s:%d", w.host, w.port)
	client, err := rpc.Dial("tcp", address)
	if err != nil {
		return err
	}
	defer func(client *rpc.Client) {
		err = client.Close()
	}(client)
	err = client.Call(s, args, result)
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

// GetAvailRegWorkers TODO: might need to return error if no workers
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

// GetRandomAvailRegWorker TODO: return error when no workers
func GetRandomAvailRegWorker() RegisteredWorker {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	index := r.Intn(len(registeredWorkers))
	return GetAvailRegWorkers()[index]
}
