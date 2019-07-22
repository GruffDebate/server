package api

import (
	_ "errors"

	"github.com/GruffDebate/server/gruff"
	arango "github.com/arangodb/go-driver"
	"github.com/labstack/echo"
)

var CTX *gruff.ServerContext
var TEST_CLIENT arango.Client
var TESTDB arango.Database

var TESTTOKEN string
var READ_ONLY bool = false

func init() {
	TEST_CLIENT, TESTDB = gruff.InitTestDB()
}

func setup() {
	//TESTDB = INITDB.Begin()

	if CTX == nil {
		CTX = &gruff.ServerContext{}
	}

	CTX.Arango.DB = TESTDB
}

func teardown() {
	//TESTDB = TESTDB.Rollback()
}

func startDBLog() {
	//TESTDB.LogMode(true)
}

func stopDBLog() {
	//TESTDB.LogMode(false)
}

func Router() *echo.Echo {
	return SetUpRouter(true, TESTDB)
}
