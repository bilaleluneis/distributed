/*	@author: Bilal El Uneis
	@since: August 2022
	@email: bilaleluneis@gmail.com	*/

package internal

import (
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	listnr, err := InitWorker(8080)
	if err != nil {
		log.Fatal("worker init failure, tests will abort")
	}
	go ProcessWorkRequest(listnr)
	os.Exit(m.Run())
}
