package db

import (
	"database/sql"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

const (
	dbName = "bookstore"
	dbUser = "root"
	dbPass = "root"
	dbURL  = "localhost"
)

var (
	dbConn *sql.DB
	onceDB sync.Once
)

func connect() *sql.DB {
	dbStr := dbUser + ":" + dbPass + "@tcp(" + dbURL + ")" + "/" + dbName
	var err error
	dbConn, err = sql.Open("mysql", dbStr)
	if err != nil {
		panic(err.Error())
	}

	return dbConn
}

func GetDBConnection() *sql.DB {
	onceDB.Do(func() {
		dbConn = connect()
	})
	err := dbConn.Ping()
	if err != nil {
		panic(err.Error())
	}

	return dbConn
}
