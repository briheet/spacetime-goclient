package main

import (
	"log"

	"github.com/briheet/spacetime-goclient/spacetimedb"
)

// For calling a reducer, first we get a dbID from a database http call, then we get token id, then we have to know before hand about the reducers that we have written in the lib.rs
// Now we make the call that we are getting after doing a tcp dump ->  sudo tcpdump -i lo port 3000 -A
func main() {

	// Need to connect to existing server running on some port
	spdb, err := spacetimedb.Connect("http://localhost", "3000", "testDB")
	if err != nil {
		log.Fatal(err)
	}

	// Defer it to close which does cleanup
	defer spdb.Disconnect()

	// First ping and check if it works or not
	err = spdb.Ping()
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Successfully connected, you are good to go!")
	}

	// Create Identity
	identity, token, err := spdb.CreateIdentity()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("identity:", identity)
	log.Println("token:", token)

	dbIden, _, _, _, err := spdb.GetDatabaseInfo("quickstart-chat")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(dbIden)

	reducerName := "send_message"
	err = spdb.SendMessageDatabase(reducerName, dbIden, token, "Hello, world")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Successfully called reducer from golang")

}
