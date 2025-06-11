package main

import (
	"log"

	"github.com/briheet/spacetimedb-go-client/spacetimedb"
)

func main() {

	spdb, err := spacetimedb.Connect("http://localhost", "3000", "testDB")
	if err != nil {
		log.Fatal(err)
	}

	defer spdb.Disconnect()


}
