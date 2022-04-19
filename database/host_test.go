package database

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestCreateHost(t *testing.T) {
	db, err := CreateInMemoryDB("../init.sql")
	if err != nil {
		t.Fatalf("Error opening database: %s", err)
	}

	id, err := db.CreateHost("example.com", "Example Server", "An example host")
	if err != nil {
		t.Errorf("Error creating host: %s", err)
	}

	host, err := db.SelectHostById(id)
	if err != nil {
		t.Errorf("Error looking up host: %s", err)
	}

	if host.ID != id {
		t.Errorf("Expected host id (%d) to match returned id (%d)", host.ID, id)
	}
	if host.Name != "Example Server" {
		t.Errorf("Expected host name to be 'Example Server' found %s", host.Name)
	}
	if len(host.Key) != TokenLength {
		t.Errorf("Expected host key to be %d bytes in length found %d", TokenLength, len(host.Key))
	}
}

var sampleHosts = []Host{
	{Domain: "google.com", Name: "Google", Description: "Google's load balancers"},
	{Domain: "google.com.au", Name: "Google AU", Description: "Google's Australian load balancers"},
	{Domain: "facebook.com", Name: "Facebook", Description: "Facebook's load balancers"},
	{Domain: "netflix.com", Name: "Netflix", Description: "Netflix's load balancers"},
	{Domain: "youtube.com", Name: "Youtube", Description: "Youtube's load balancers"},
	{Domain: "news.ycombinator.com", Name: "Hacker News", Description: "Ycombinator's tech news"},
	{Domain: "reddit.com", Name: "Reddit", Description: "Reddit's load balancers"},
}

func TestHostSearchExact(t *testing.T) {
	db, err := CreateInMemoryDB("../init.sql")
	if err != nil {
		t.Errorf("Error opening database: %s", err)
	}

	for _, host := range sampleHosts {
		_, err := db.CreateHost(host.Domain, host.Name, host.Description)
		if err != nil {
			t.Errorf("Error creating host: %s", err)
		}
	}

	hosts, err := db.SearchHostByDomain("news.ycombinator.com")
	if err != nil {
		t.Errorf("Error searching host: %s", err)
	}

	if len(hosts) != 1 {
		t.Errorf("Expected 1 result found %d", len(hosts))
	}
	if hosts[0].Name != "Hacker News" {
		t.Errorf("Expected host name to be 'Hacker News' found %s", hosts[0].Name)
	}
}

func TestHostSearchPartial(t *testing.T) {
	db, err := CreateInMemoryDB("../init.sql")
	if err != nil {
		t.Errorf("Error opening database: %s", err)
	}

	for _, host := range sampleHosts {
		_, err := db.CreateHost(host.Domain, host.Name, host.Description)
		if err != nil {
			t.Errorf("Error creating host: %s", err)
		}
	}

	hosts, err := db.SearchHostByDomain("google")
	if err != nil {
		t.Errorf("Error searching host: %s", err)
	}

	if len(hosts) != 2 {
		t.Errorf("Expected 2 result found %d", len(hosts))
	}
	if hosts[0].Name != "Google" {
		t.Errorf("Expected host name to be 'Google' found %s", hosts[0].Name)
	}
	if hosts[1].Name != "Google AU" {
		t.Errorf("Expected host name to be 'Google AU' found %s", hosts[0].Name)
	}
}
