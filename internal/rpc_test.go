package internal

import (
	"distributed/common"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strings"
	"testing"
)

func TestRpcNode(t *testing.T) {
	var uuid common.UUID

	// insert call test
	err := client.Call(INSERT, []byte{}, &uuid)
	if err != nil || uuid == "" {
		t.Fail()
	}

	//retrieve call with remove
	location := common.Location{Uuid: uuid}
	params := common.SearchParams{Address: location, Remove: true}
	var result []byte
	err = client.Call(RETRIEVE, params, &result)
	if err != nil {
		t.Fail()
	}

	// retrieve call : check if entry was removed
	err = client.Call(RETRIEVE, params, &result)
	if !strings.Contains(fmt.Sprint(err), "nothing to pop") {
		t.Fail()
	}
}

var client *rpc.Client

func TestMain(m *testing.M) {

	go func() {
		err := InitServer(8080)
		if err != nil {
			log.Fatal("test server startup error", err)
		}
	}()

	var err error
	client, err = rpc.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal("test client init failed", err)
	}

	os.Exit(m.Run())

}
