package config

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var CONFIGURATIONS map[string]string = map[string]string{
	"GRUFF_ENV":  "local",
	"GRUFF_DB":   "host=gruff.c7qnzdzjyjrm.us-west-2.rds.amazonaws.com user=gruff dbname=gruff password=gruffdeveloper7240 sslmode=disable",
	"GRUFF_NAME": "GRUFF",
}

func Init() {
	if os.Getenv("GRUFF_NAME") == "" {
		os.Setenv("GRUFF_NAME", CONFIGURATIONS["GRUFF_NAME"])
	}
}

func InitDB() (rw *gorm.DB) {
	if os.Getenv("GRUFF_DB") == "" {
		os.Setenv("GRUFF_DB", CONFIGURATIONS["GRUFF_DB"])
	}

	db, err := gorm.Open("postgres", os.Getenv("GRUFF_DB"))
	if err != nil {
		panic(err.Error())
	}
	rw = db
	fmt.Println("Initialized read-write database connection pool")

	return
}
