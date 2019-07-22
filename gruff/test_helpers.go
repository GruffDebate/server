package gruff

import (
	"fmt"
	"os"

	arango "github.com/arangodb/go-driver"
)

func InitTestDB() (arango.Client, arango.Database) {
	if os.Getenv("ARANGO_ENDPOINT") == "" {
		os.Setenv("ARANGO_ENDPOINT", "http://localhost:8529")
	}
	if os.Getenv("ARANGO_DB") == "" {
		os.Setenv("ARANGO_DB", "gruff_test")
	}
	if os.Getenv("ARANGO_USER") == "" {
		os.Setenv("ARANGO_USER", "root")
	}
	if os.Getenv("ARANGO_PASS") == "" {
		os.Setenv("ARANGO_PASS", "")
	}

	client, err := OpenTestConnection()
	if err != nil {
		fmt.Println("No error should happen when connecting to test database, but got:", err)
	}

	db, err := OpenArangoDatabase(client)
	if err != nil {
		fmt.Println("Error opening the test database:", err)
	}

	cleanData(db)

	return client, db
}

func OpenTestConnection() (arango.Client, error) {
	return OpenArangoConnection()
}

func cleanData(db arango.Database) {
	ctx := ArangoContext{DB: db}

	models := []ArangoObject{
		User{},
		Inference{},
		BaseClaimEdge{},
		PremiseEdge{},
		Argument{},
		Claim{},
	}

	for _, m := range models {
		col, err := ctx.CollectionFor(m)
		if err != nil {
			// bummer
			fmt.Println("Error getting collection for model")
		}
		if err := col.Truncate(nil); err != nil {
			// bummer
			fmt.Println("Truncating collection")
		}
	}

}
