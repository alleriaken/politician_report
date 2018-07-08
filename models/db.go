package models

import (
	"database/sql"
	"log"
	"os"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitDB() {
	var err error

	SQL_HOST := os.Getenv("SQL_HOST")
	SQL_USERNAME := os.Getenv("SQL_USERNAME")
	SQL_PASSWORD := os.Getenv("SQL_PASSWORD")
	//SQL_PORT := os.Getenv("SQL_PORT")
	SQL_DBNAME := os.Getenv("SQL_DBNAME")

	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@%s/%s", SQL_USERNAME, SQL_PASSWORD, SQL_HOST, SQL_DBNAME))

	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	if err != nil {
		log.Panic(err)
	}

	if err = db.Ping(); err != nil {
		log.Panic(err)
	}
}

func CloseDB()  {
	db.Close()
}