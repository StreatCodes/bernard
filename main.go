package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/StreatCodes/bernard/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
)

const DBFile = "./bernard.db"

type Service struct {
	DB *database.DB
}

func main() {
	//Check DB is setup
	_, err := os.Stat(DBFile)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Printf("%s not found, running setup\n", DBFile)
		err := database.SetupDB(DBFile, "init.sql")
		if err != nil {
			log.Fatalf("Error initializing DB: %s\n", err)
		}

		fmt.Printf("Database setup complete! You can always delete %s to start fresh\n", DBFile)
	} else if err != nil {
		log.Fatalf("Error stating DB file (%s): %s", DBFile, err)
	}

	//Connect to DB
	fmt.Printf("Connecting to %s\n", DBFile)
	db, err := database.ConnectFileDB(DBFile)
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

	r.With(service.hostAuthMiddleware).Post("/log", service.HandleNewLog)

	r.Route("/", func(r chi.Router) {
		r.Use(service.userAuthMiddleWare)

		r.Post("/host", service.HandleCreateHost)
	})

	//Run web server
	webAddr := "localhost:3000"
	fmt.Printf("Starting web server on %s\n", webAddr)
	http.ListenAndServe(webAddr, r)
}
