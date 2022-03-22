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

type Service struct {
	DB *database.DB
}

func main() {
	//Setup DB
	db, err := database.CreateInMemoryDB("init.sql")
	if err != nil {
		log.Fatalf("Error opening database: %s\n", err)
	}

	service := Service{
		DB: db,
	}

	//Setup middlewares
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	//Setup routes
	r.Post("/login", service.HandleLogin)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	http.ListenAndServe(":3000", r)
}
