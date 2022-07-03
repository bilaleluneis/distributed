/*	@author: Bilal El Uneis
	@since: June 2022
	@email: bilaleluneis@gmail.com	*/

package distributed

import (
	"testing"
)

func TestTcpRpc(t *testing.T) {

	go func() {
		_ = AsServer(8080)
	}()

	_ = AsClient("localhost", 8080)

	node := Node{}

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
