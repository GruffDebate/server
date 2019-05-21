package gruff

import (
	"context"
	"fmt"
	"os"

	arango "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
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

	ctx := context.Background()
	db, err := client.Database(ctx, os.Getenv("ARANGO_DB"))
	if err != nil {
		fmt.Println("Error opening the test database:", err)
	}

	cleanData(db)

	return client, db
}

func OpenTestConnection() (arango.Client, error) {
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{os.Getenv("ARANGO_ENDPOINT")},
	})
	if err != nil {
		return nil, err
	}
	conn, err = conn.SetAuthentication(arango.BasicAuthentication(os.Getenv("ARANGO_USER"), os.Getenv("ARANGO_PASS")))
	if err != nil {
		return nil, err
	}
	db, err := arango.NewClient(arango.ClientConfig{
		Connection: conn,
	})

	return db, err
}

func cleanData(db arango.Database) {
	ctx := ArangoContext{DB: db}

	models := []ArangoObject{
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
