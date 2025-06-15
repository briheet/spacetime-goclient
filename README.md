# Spacetime GoClient

## _Go client for SpaceTimeDB_

This repository includes go client for SpacetimeDB. Currently under development.

Currently includes:

- Connect to required server and Ping to check
- Identity API integration
- Database API integration

Start a new server via docker this:

```bash
docker run --rm --pull always -p 3000:3000 clockworklabs/spacetime start
```

This is a external package, make sure golang is installed on your system. Then, install it via this command:

```go
go get -u github.com/briheet/spacetime-goclient/spacetimedb
```

To start using, follow this:

```go
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

}
```

## Identity

1. To create a spacetime public identities and private tokens:

```go
	// Create Identity and Token
	identity, token, err := spdb.CreateIdentity()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("identity:", identity)
	log.Println("token:", token)
```

2. To generate short-lived access token which can be used in untrusted contexts:

```go
	// Create Websocket token
	websocketToken, err := spdb.CreateIdentityWebsocketToken()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("websocketToken:", websocketToken)
```

3. To fetch the public key used by the database to verify tokens

```go
	// Get public key used by the database to verify tokens
	publicKey, err := spdb.GetPublicKey()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("publicKey:", publicKey)
```

4. Associate an email with a Spacetime identity (Needs some fixes)

```go
	// Register identity with email. Currently endpoint issue
	emailIdentity, emailToken, err := spdb.RegisterIdentityWithEmail("briheetyadav@gmail.com")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Email identity:", emailIdentity)
	log.Println("Email token:", emailToken)
```

5. To list all databases owned by an identity

```go
	// Create Identity and Token
	identity, _, err := spdb.CreateIdentity()
	if err != nil {
		log.Fatal(err)
	}

	// Get databases by identity
	databases, err := spdb.GetDatabasesByIdentity(identity)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(databases)
```

6. To verify identity and token is valid or not

```go
	// Create Identity and Token
	identity, token, err := spdb.CreateIdentity()
	if err != nil {
		log.Fatal(err)
	}

	// Verify identity and token
	if err := spdb.VerifyIdentityToken(identity, token); err != nil {
		log.Fatalf("Identity verification failed: %v", err)
	} else {
		log.Println("✅ Identity and token verified successfully.")
	}
```

## Database

1. Get a database's identity, owner identity, host type, number of replicas and a hash of its WASM module.

```go
	dbIden, _, _, _, err := spdb.GetDatabaseInfo("quickstart-chat")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(dbIden)
```

2. Delete a database.

```go
	err = spdb.DeleteDatabase(dbIden, token)
	if err != nil {
		log.Fatal(err)
	}
```

3. Get the names this database can be identified by.

```go
	names, err := spdb.GetDatabaseNames(dbIden)
	if err != nil {
		log.Fatalf("Error fetching names: %v", err)
	}
	log.Println("Names:", names)
```

4. Get the names this database can be identified by.

```go
	names, err := spdb.GetDatabaseNames(dbIden)
	if err != nil {
		log.Fatalf("Error fetching names: %v", err)
	}
	log.Println("Names:", names)
```

5. Add a new name for this database.

```go
	err = spdb.AddDatabaseName(dbIden, "mychat", token)
	if err != nil {
		log.Fatalf("Failed to name database: %v", err)
	}
	log.Println("Database name assigned successfully.")
```

6. Get the identity of a database.

```go
	id, err := spdb.GetDatabaseIdentity("quickstart-chat")
	if err != nil {
		log.Fatalf("Failed to fetch identity: %v", err)
	}
	log.Println("Database identity:", id)
```

7. Begin a WebSocket connection with a database.

```go
	conn, err := spdb.WebsocketSubscribe("quickstart-chat", token, "")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket read error:", err)
			break
		}
		log.Println("Received message:", string(msg))
	}
```

8. Invoke a reducer in a database.

```go
	reducerName := "send_message"
	err = spdb.SendMessageDatabase(reducerName, dbIden, token, "Hello, world")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Successfully called reducer from golang")
```

9. Retrieve logs from a database.

```go
	logs, err := spdb.GetDatabaseLogs("quickstart-chat", token, 100, false)
	if err != nil {
		log.Fatalf("failed to get logs: %v", err)
	}
	defer logs.Close()

	scanner := bufio.NewScanner(logs)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
```

10. Run a SQL query against a database.

```go
	dbName := "quickstart-chat"
	results, err := spdb.RunSQLQuery(`SELECT * FROM person;`, token, dbName)
	if err != nil {
		log.Fatal(err)
	}

	for _, res := range results {
		fmt.Printf("Schema: %+v\n", res.Schema)
		fmt.Printf("Rows: %+v\n", res.Rows)
	}
```
