package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
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

func initDB(test bool) (*database.DB, error) {
	if test {
		return database.CreateInMemoryDB("./init.sql")
	}

	//Check DB is setup
	_, err := os.Stat(DBFile)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Printf("%s not found, running setup\n", DBFile)
		err := database.SetupDB(DBFile, "init.sql")
		if err != nil {
			return nil, err
		}

		fmt.Printf("Database setup complete! You can always delete %s to start fresh\n", DBFile)
	} else if err != nil {
		return nil, err
	}

	//Connect to DB
	fmt.Printf("Connecting to %s\n", DBFile)
	return database.ConnectFileDB(DBFile)
}

func (s *Service) startWebServer(addr string) error {
	//Setup middlewares
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))
	r.Route("/api", func(r chi.Router) {

		//Setup routes
		r.Post("/login", s.HandleLogin)

		r.With(s.hostAuthMiddleware).Post("/log", s.HandleNewLog)

		r.Route("/", func(r chi.Router) {
			r.Use(s.userAuthMiddleWare)

			r.Post("/host", s.HandleCreateHost)
			r.Get("/host", s.HandleGetHosts)
			r.Get("/host/{id}", s.HandleGetHost)
			r.Get("/host/{id}/log", s.HandleGetHostLogs)
			// r.Get("/host/{id}/metrics", s.HandleGetHostMetrics)
		})
	})

	r.NotFound(fileHandler)

	//Run web server
	fmt.Printf("Starting web server on %s\n", addr)
	return http.ListenAndServe(addr, r)
}

func fileHandler(w http.ResponseWriter, r *http.Request) {
	reqPath := path.Join("frontend", path.Clean(r.URL.Path))

	//serve index.html if requested file doesn't exist
	_, err := os.Stat(reqPath)
	if errors.Is(err, os.ErrNotExist) {
		http.ServeFile(w, r, "frontend/index.html")
		return
	}

	http.ServeFile(w, r, reqPath)
}

func main() {
	db, err := initDB(false)
	if err != nil {
		log.Fatalf("Error intilizing DB: %s\n", err)
	}

	service := Service{
		DB: db,
	}

	err = service.startWebServer("localhost:3000")
	if err != nil {
		log.Fatalf("Error starting web server: %s\n", err)
	}
}
