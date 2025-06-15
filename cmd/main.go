package main

import (
	"fmt"
	"log"

	"github.com/briheet/spacetime-goclient/spacetimedb"
)

// For calling a reducer, first we get a dbID from a database http call, then we get token id, then we have to know before hand about the reducers that we have written in the lib.rs
// Now we make the call that we are getting after doing a tcp dump ->  sudo tcpdump -i lo port 3000 -A
func main() {

	// Need to connect to existing server running on some port
	spdb, err := spacetimedb.Connect("localhost", "3000", "testDB")
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

	// reducerName := "send_message"
	// err = spdb.SendMessageDatabase(reducerName, dbIden, token, "Hello, world")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// log.Println("Successfully called reducer from golang")

	// err = spdb.DeleteDatabase(dbIden, token)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// 5.
	// names, err := spdb.GetDatabaseNames(dbIden)
	// if err != nil {
	// 	log.Fatalf("Error fetching names: %v", err)
	// }
	// log.Println("Names:", names)

	// 6.
	// err = spdb.AddDatabaseName(dbIden, "mychat", token)
	// if err != nil {
	// 	log.Fatalf("Failed to name database: %v", err)
	// }
	// log.Println("Database name assigned successfully.")

	// 7.
	// id, err := spdb.GetDatabaseIdentity("quickstart-chat")
	// if err != nil {
	// 	log.Fatalf("Failed to fetch identity: %v", err)
	// }
	// log.Println("Database identity:", id)

	// 8. Websocket
	// conn, err := spdb.WebsocketSubscribe("quickstart-chat", token, "")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer conn.Close()
	//
	// for {
	// 	_, msg, err := conn.ReadMessage()
	// 	if err != nil {
	// 		log.Println("WebSocket read error:", err)
	// 		break
	// 	}
	// 	log.Println("Received message:", string(msg))
	// }

	// 9.
	// logs, err := spdb.GetDatabaseLogs("quickstart-chat", token, 100, false)
	// if err != nil {
	// 	log.Fatalf("failed to get logs: %v", err)
	// }
	// defer logs.Close()
	//
	// scanner := bufio.NewScanner(logs)
	// for scanner.Scan() {
	// 	fmt.Println(scanner.Text())
	// }

	// 10.
	dbName := "quickstart-chat"
	results, err := spdb.RunSQLQuery(`SELECT * FROM person;`, token, dbName)
	if err != nil {
		log.Fatal(err)
	}

	for _, res := range results {
		fmt.Printf("Schema: %+v\n", res.Schema)
		fmt.Printf("Rows: %+v\n", res.Rows)
	}

}
