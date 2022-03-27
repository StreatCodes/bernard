package database

import (
	"log"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func TestLogsBetween(t *testing.T) {
	db, err := CreateInMemoryDB("../init.sql")
	if err != nil {
		t.Errorf("Error opening database: %s", err)
	}

	logData := []Log{
		{Time: time.Now(), HostID: 1, Level: 0, Service: "test", Content: "Logged now"},
		{Time: time.Now().Add(30 * time.Minute), HostID: 1, Level: 0, Service: "test", Content: "Logged 30 minutes ago"},
		{Time: time.Now().Add(-30 * time.Minute), HostID: 1, Level: 0, Service: "test", Content: "Logged 30 minutes from now"},
		{Time: time.Now().Add(-time.Hour), HostID: 1, Level: 0, Service: "test", Content: "Logged 1 hour ago"},
	}

	//fill DB with sample logs
	for _, data := range logData {
		err = db.InsertLog(data)
		if err != nil {
			t.Errorf("Error inserting log: %s", err)
		}
	}

	logs, err := db.SelectLogsBetween(time.Now().Add(-5*time.Minute), time.Now().Add(5*time.Minute))
	if err != nil {
		log.Fatalf("Error fetching logs: %s", err)
	}

	logCount := len(logs)
	if logCount != 1 {
		t.Errorf("Expected 1 logs added, found %d", logCount)
	}
	logContent := logs[0].Content
	if logContent != "Logged now" {
		t.Errorf("Expected log content to contain 'Logged now', found %s", logContent)
	}
}

func TestHostLogsBetween(t *testing.T) {
	db, err := CreateInMemoryDB("../init.sql")
	if err != nil {
		t.Errorf("Error opening database: %s", err)
	}

	logData := []Log{
		{Time: time.Now(), HostID: 1, Level: 0, Service: "host 1", Content: "Logged now on host 1"},
		{Time: time.Now().Add(time.Second), HostID: 2, Level: 0, Service: "host 2", Content: "Logged now on host 2"},
		{Time: time.Now().Add(2 * time.Second), HostID: 3, Level: 0, Service: "host 3", Content: "Logged now on host 3"},
	}

	//fill DB with sample logs
	for _, data := range logData {
		err = db.InsertLog(data)
		if err != nil {
			t.Errorf("Error inserting log: %s", err)
		}
	}

	logs, err := db.SelectHostLogsBetween(int64(2), time.Now().Add(-10*time.Minute), time.Now().Add(10*time.Minute))
	if err != nil {
		log.Fatalf("Error fetching logs: %s", err)
	}

	logCount := len(logs)
	if logCount != 1 {
		t.Errorf("Expected 2 logs added, found %d", logCount)
	}
	if logs[0].Service != "host 2" {
		t.Errorf("Expected service name 'host 2' found '%s'", logs[0].Service)
	}
}
