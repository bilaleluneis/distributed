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

// Worker TODO: implement asyncInvoke
type Worker struct {
	Host string
	Port int
}

func (w Worker) Invoke(s Service, args any, result any) error {
	address := fmt.Sprintf("%s:%d", w.Host, w.Port)
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

var workers []Worker

// Init performs initalization of workers once
func Init(w []Worker) {
	if len(workers) == 0 {
		workers = w
	}
}

// Get this will return workers in random order
func Get() []Worker {
	copyOfWorkers := make([]Worker, 0)
	for _, w := range workers {
		copyOfWorkers = append(copyOfWorkers, w)
	}
	rand.Seed(time.Now().Unix())
	rand.Shuffle(len(copyOfWorkers), func(i, j int) {
		copyOfWorkers[i] = copyOfWorkers[j]
		copyOfWorkers[j] = copyOfWorkers[i]
	})
	return copyOfWorkers
}
