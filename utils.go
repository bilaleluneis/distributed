/*	@author: Bilal El Uneis
	@since: September 2022
	@email: bilaleluneis@gmail.com	*/

package distributed

import (
	"distributed/common"
	"distributed/internal"
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"os"
)

type Connection struct {
	Host string
	Port int
}

func isWorker() (port int, isWorker bool) {
	if len(os.Args) > 1 {
		flag.IntVar(&port, "p", 8080, "Specify Port, Default is 8080")
		flag.Usage = func() {
			fmt.Printf("to start as worker pass -p port number ex: -p 8090\n")
			fmt.Printf("otherwise omit flags to start as client\n")
		}
		flag.Parse()
	}
	isWorker = port != 0
	return
}

func Init(withWorkers ...Connection) {
	//FIXME: take this out once you refactor how functional types are registered with GOB
	gob.Register(internal.Compute{})

	for _, worker := range withWorkers {
		common.RegisterWorker(worker.Host, worker.Port)
	}
	if port, ok := isWorker(); ok {
		var worker internal.Worker
		var err error
		if worker, err = internal.NewWorker(port); err != nil {
			log.Panicf("failed to create worker with error %s", err.Error())
		}
		worker.Start()
	}
}

func RegisterFilter[T any](f ...common.Filterer[T]) {
	for _, filter := range f {
		gob.Register(filter)
		gob.Register(internal.Filter[T]{WithFilter: filter})
	}
}

func RegisterMapper[T any, R any](m ...common.Mapper[T, R]) {
	for _, mapper := range m {
		gob.Register(mapper)
		gob.Register(internal.Map[T, R]{WithMapper: mapper})
	}
}

func RegisterReducer[T any, R any](r ...common.Reducer[T, R]) {
	for _, reducer := range r {
		gob.Register(reducer)
		gob.Register(internal.Reduce[T, R]{WithReducer: reducer})
	}
}
