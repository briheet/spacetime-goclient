# Spacetime GoClient

## _Go client for SpaceTimeDB_

This repository includes go client for SpacetimeDB. Currently under development.

Currently includes:
- Connect to required server
- Ping it to check if connected or not 

Start a new server via docker this:
```bash
docker run --rm --pull always -p 3000:3000 clockworklabs/spacetime start
```


This is a external package, make sure golang is installed on your system. Then,  install it via this command:
```go
go get -u github.com/briheet/spacetime-goclient
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