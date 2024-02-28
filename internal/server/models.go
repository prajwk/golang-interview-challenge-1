package server

import "time"

type StockInfo struct {
	Ticker     string    `json:"ticker" db:"ticker"`
	Timestamp  time.Time `json:"timestamp" db:"timestamp"`
	Open       float64   `json:"open" db:"open"`
	High       float64   `json:"high" db:"high"`
	Low        float64   `json:"low" db:"low"`
	Close      float64   `json:"close" db:"close"`
	Volume     int64     `json:"volume" db:"volume"`
	InsertedDt time.Time `json:"insertedDt" db:"inserted_dt"`
}
