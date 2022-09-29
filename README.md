# distributed: API That Simplifies the Creation of Distributed Data Structures

distributed is generic, high level RPC based api that allows the creation of data structures that live on cluster. 

Here's an example of a distributed list used in local cluster, to start as worker , run the program with -p 8081.
to start as client run without options:

```go
package main

import (
	"distributed"
	"distributed/common"
	"distributed/datastructure"
	"log"
	"math/rand"
	"time"
)

type lessThanfilter struct{ Value int }

func (f lessThanfilter) Filter(n common.NodeLike[int]) bool {
	return n.GetData() < f.Value
}

func main() {
	
	distributed.RegisterFilter[int](lessThanfilter{})
	distributed.Init(distributed.Connection{Host: "localhost", Port: 8081})

	rand.Seed(time.Now().Unix())
	l, err := datastructure.NewList[int](rand.Perm(100)...)
	if err != nil {
		log.Panicf("client side error due %s", err.Error())
	}

	log.Printf("list size is %d", l.Len())
	log.Printf("list group ID is %s", l.Identity())
	log.Printf("Performing Filter/OP")

	err = datastructure.Filter[int](l, lessThanfilter{10})
	if err != nil {
		log.Panicf("filter op failed due %s", err.Error())
	}

	var result []int
	if result, err = datastructure.Compute[int](l); err != nil {
		log.Panicf("compute op failed due %s", err.Error())
	}

	log.Printf("Compute Result is %d", result)

}
```
## License
distributed is licensed under the [MIT License](LICENSE).