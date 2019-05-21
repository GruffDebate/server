package gruff

import (
	arango "github.com/arangodb/go-driver"
)

var TEST_CLIENT arango.Client
var TESTDB arango.Database
var CTX *ServerContext

func init() {
	CTX = &ServerContext{}
	TEST_CLIENT, TESTDB = InitTestDB()
}

func setupDB() {
	CTX.Arango.DB = TESTDB
}

func teardownDB() {
}

func startDBLog() {
}

func stopDBLog() {
}
