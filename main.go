package main

import (
	"log"
	"net/http"
	"time"

	"github.com/StreatCodes/bernard/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
)

type Log struct {
	At      time.Time
	Content string
}

func main() {
	_, err := database.CreateInMemoryDB("init.sql")
	if err != nil {
		log.Fatalf("Error opening database: %s\n", err)
	}

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	http.ListenAndServe(":3000", r)
}
