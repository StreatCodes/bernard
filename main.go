package main

import (
	"fmt"
	"log"
	"time"

	"github.com/StreatCodes/bernard/database"
	_ "github.com/mattn/go-sqlite3"
)

type Log struct {
	At      time.Time
	Content string
}

func main() {
	db, err := database.CreateInMemoryDB("init.sql")
	if err != nil {
		log.Fatalf("Error opening database: %s\n", err)
	}

	err = db.InsertLog(database.Log{
		Time:    time.Now(),
		HostID:  1,
		Level:   0,
		Service: "test",
		Content: "Testing 123",
	})
	if err != nil {
		log.Fatalf("Error inserting log: %s\n", err)
	}

	logs, err := db.SelectSince(time.Now().Add(-time.Hour))
	if err != nil {
		log.Fatalf("Error fetching logs: %s\n", err)
	}

	for _, log := range logs {
		fmt.Printf("%s\n", log.Time)
	}
}
