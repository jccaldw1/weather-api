package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("hello world")
	connStr := os.Getenv("CONNECTION_STRING")
	if connStr == "" {
		log.Fatal("Could not find connection string")
		return
	}

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal("Could not open db connection: ", err)
		return
	}

	defer db.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/", helloWorldHandler)
	mux.HandleFunc("/date/{date}", GetDateHandler)

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal("Could not start server.")
	}
}

func GetDateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Getting records for the passed date (year, month, day)")
	date := r.PathValue("date")
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HelloWorld")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := struct {
		Key string `json:"keyy"`
	}{
		Key: "hello",
	}

	json.NewEncoder(w).Encode(response)
}
