package gruff

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var INITDB *gorm.DB
var TESTDB *gorm.DB
var CTX ServerContext

func init() {
	INITDB = InitTestDB()
}

func setupDB() {
	TESTDB = INITDB.Begin()
	CTX.Database = TESTDB
}

func teardownDB() {
	TESTDB = TESTDB.Rollback()
}

func startDBLog() {
	TESTDB.LogMode(true)
}

func stopDBLog() {
	TESTDB.LogMode(false)
}
