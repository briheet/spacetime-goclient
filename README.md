# Spacetime GoClient

## _Go client for SpaceTimeDB_

This repository includes go client for SpacetimeDB. Currently under development.

Currently includes:

- Connect to required server
- Ping it to check if connected or not
- Create spacetime public identities and private tokens.

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
	// Get databases by identity
	databases, err := spdb.GetDatabasesByIdentity(identity)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(databases)
```

6. To verify identity and token is valid or not

```go
	// Verify identity and token
	if err := spdb.VerifyIdentityToken(identity, token); err != nil {
		log.Fatalf("Identity verification failed: %v", err)
	} else {
		log.Println("✅ Identity and token verified successfully.")
	}
```
