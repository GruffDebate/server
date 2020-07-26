package config

import (
	"fmt"
	"os"

	"github.com/GruffDebate/server/gruff"
	arango "github.com/arangodb/go-driver"
)

var CONFIGURATIONS map[string]string = map[string]string{
	"GRUFF_ENV":            "development",
	"GRUFF_NAME":           "GRUFF",
	"PORT":                 "8080",
	"ARANGO_ENDPOINT":      "http://localhost:8529",
	"ARANGO_DB":            "gruff",
	"ARANGO_USER":          "root",
	"ARANGO_PASS":          "",
	"JWT_KEY_SIGNIN":       "a324dd15-74c5-44ea-8f64-8f0e6b90844c",
	"JWT_TOKEN_EXPIRATION": "720",
}

func Init() {
	if os.Getenv("GRUFF_ENV") == "" {
		os.Setenv("GRUFF_ENV", CONFIGURATIONS["GRUFF_ENV"])
	}
	if os.Getenv("GRUFF_NAME") == "" {
		os.Setenv("GRUFF_NAME", CONFIGURATIONS["GRUFF_NAME"])
	}
	if os.Getenv("PORT") == "" {
		os.Setenv("PORT", CONFIGURATIONS["PORT"])
	}
	if os.Getenv("JWT_KEY_SIGNIN") == "" {
		os.Setenv("JWT_KEY_SIGNIN", CONFIGURATIONS["JWT_KEY_SIGNIN"])
	}
	if os.Getenv("JWT_TOKEN_EXPIRATION") == "" {
		os.Setenv("JWT_TOKEN_EXPIRATION", CONFIGURATIONS["JWT_TOKEN_EXPIRATION"])
	}
	if os.Getenv("ARANGO_ENDPOINT") == "" {
		os.Setenv("ARANGO_ENDPOINT", CONFIGURATIONS["ARANGO_ENDPOINT"])
	}
	if os.Getenv("ARANGO_DB") == "" {
		os.Setenv("ARANGO_DB", CONFIGURATIONS["ARANGO_DB"])
	}
	if os.Getenv("ARANGO_USER") == "" {
		os.Setenv("ARANGO_USER", CONFIGURATIONS["ARANGO_USER"])
	}
	if os.Getenv("ARANGO_PASS") == "" {
		os.Setenv("ARANGO_PASS", CONFIGURATIONS["ARANGO_PASS"])
	}

	fmt.Println("GRUFF_ENV=", os.Getenv("GRUFF_ENV"))
	fmt.Println("GRUFF_NAME=", os.Getenv("GRUFF_NAME"))
	fmt.Println("PORT=", os.Getenv("PORT"))
	fmt.Println("JWT_KEY_SIGNIN=", os.Getenv("JWT_KEY_SIGNIN"))
	fmt.Println("JWT_TOKEN_EXPIRATION=", os.Getenv("JWT_TOKEN_EXPIRATION"))
	fmt.Println("ARANGO_ENDPOINT=", os.Getenv("ARANGO_ENDPOINT"))
	fmt.Println("ARANGO_DB=", os.Getenv("ARANGO_DB"))
	fmt.Println("ARANGO_USER=", os.Getenv("ARANGO_USER"))
	fmt.Println("ARANGO_PASS=", os.Getenv("ARANGO_PASS"))
}

func InitDB() arango.Database {
	client, err := gruff.OpenTestConnection()
	if err != nil {
		fmt.Println("No error should happen when connecting to test database, but got:", err)
	}

	db, err := gruff.OpenArangoDatabase(client)
	if err != nil {
		fmt.Println("Error opening the test database:", err)
	}

	return db
}
