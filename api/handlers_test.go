package api

import (
	_ "errors"

	"github.com/GruffDebate/server/gruff"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo"
)

var CTX *gruff.ServerContext
var INITDB *gorm.DB
var TESTDB *gorm.DB

var TESTTOKEN string
var READ_ONLY bool = false

func init() {
	INITDB = gruff.InitTestDB()
}

func setup() {
	TESTDB = INITDB.Begin()

	if CTX == nil {
		CTX = &gruff.ServerContext{}
	}

	CTX.Database = TESTDB
}

func teardown() {
	TESTDB = TESTDB.Rollback()
}

func startDBLog() {
	TESTDB.LogMode(true)
}

func stopDBLog() {
	TESTDB.LogMode(false)
}

func Router() *echo.Echo {
	return SetUpRouter(true, TESTDB)
}
