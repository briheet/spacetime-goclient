package main

import (
	"log"

	"github.com/briheet/spacetime-goclient/spacetimedb"
)

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

	// Create Websocket token
	websocketToken, err := spdb.CreateIdentityWebsocketToken()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("websocketToken:", websocketToken)

	// Get public key used by the database to verify tokens
	publicKey, err := spdb.GetPublicKey()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("publicKey:", publicKey)
}
