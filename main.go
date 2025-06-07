package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("hello world")
	connStr := os.Getenv("CONNECTION_STRING")
	if connStr == "" {
		log.Fatal("Could not open connection")
		return
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Could not open db connection: ", err)
		return
	}
	defer db.Close()
}
