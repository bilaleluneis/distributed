/*	@author: Bilal El Uneis
	@since: June 2022
	@email: bilaleluneis@gmail.com	*/

package distributed

import (
	"log"
	"os"
	"testing"
)

func TestTcpRpc(t *testing.T) {

	node := Node[int]{}

	err := node.Push(1)

	if err != nil {
		t.Errorf("node.push() call failed")
	}

	//do a pop call to see if you get value back
	result, err := node.Pop(false)

	if err != nil {
		t.Errorf("node.pop() call failed")
	}

	if result != 1 {
		t.Errorf("node.pop() result is wrong")
	}

}

// allow single initialization of client and server instance for all tests
func TestMain(m *testing.M) {

	go func() {
		err := AsServer(8080)
		if err != nil {
			log.Fatal("test server startup error", err)
		}
	}()

	err := AsClient("localhost", 8080)
	if err != nil {
		log.Fatal("error aquiring client instance", err)
	}

	exitVal := m.Run()
	os.Exit(exitVal)
}
