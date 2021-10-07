package database

import (
	"database/sql"
	"fmt"
	"log"
)

const (
	username = "root"
	password = ""
	hostname = "127.0.0.1:3306"
	dbname   = "go_crawler"
)

func DBConn() *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbname))
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
	}
	return db
}
