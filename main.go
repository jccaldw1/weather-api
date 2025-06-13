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

var globalDb *sql.DB

func main() {
	connStr := os.Getenv("CONNECTION_STRING")
	if connStr == "" {
		log.Fatal("Could not find connection string")
		return
	}

	db, err := sql.Open("postgres", connStr)
	globalDb = db

	if err != nil {
		log.Fatal("Could not open db connection: ", err)
		return
	}

	defer db.Close()

	db.Ping()
	fmt.Println("after ping 1")

	mux := http.NewServeMux()
	mux.HandleFunc("/date/{date}", GetDateHandler)

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal("Could not start server.")
	}
}

type Response struct {
	Date         string  `json:"date"`
	High         float32 `json:"high"`
	Low          float32 `json:"low"`
	DaysAhead    int     `json:"days_ahead"`
	DateRecorded string  `json:"date_recorded"`
}

func GetDateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Getting records for the passed date (year, month, day)")
	date := r.PathValue("date")
	query := fmt.Sprintf("SELECT * FROM \"WeatherRecord\" WHERE date = '%s'", date)
	rows, err := globalDb.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var jsonResponse []Response

	for rows.Next() {
		var id int
		var date string
		var high float32
		var low float32
		var daysAhead int
		var dateRecorded string

		err := rows.Scan(&id, &date, &high, &low, &daysAhead, &dateRecorded)
		if err != nil {
			log.Fatal("Could not scan row: ", err)
		}
		jsonResponse = append(jsonResponse, Response{date, high, low, daysAhead, dateRecorded})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(jsonResponse)

	// Now that we've sent back the response, we should also send back metrics to the database (the metrics are lazily evaluated).

	haveDataForToday := false
	var todayActualHigh float32
	var todayActualLow float32

	for _, response := range jsonResponse {
		if response.DaysAhead == 0 {
			haveDataForToday = true
			todayActualHigh = response.High
			todayActualLow = response.Low
			break
		}
	}

	if haveDataForToday {
		type Margin struct {
			DateRecorded string  `json:"date_recorded"`
			Date         string  `json:"date"`
			DaysAhead    int     `json:"days_ahead"`
			HighMargin   float32 `json:"high_margin"`
			LowMargin    float32 `json:"low_margin"`
		}

		var margins []Margin

		for _, response := range jsonResponse {
			margins = append(margins, Margin{
				DateRecorded: response.DateRecorded,
				Date:         response.Date,
				DaysAhead:    response.DaysAhead,
				HighMargin:   (response.High - todayActualHigh) / todayActualHigh,
				LowMargin:    (response.Low - todayActualLow) / todayActualLow,
			})
		}

		for _, margin := range margins {

		}
	}
}
