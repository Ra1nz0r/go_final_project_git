package server

import (
	"log"
	"net/http"
	"time"

	"github.com/ra1nz0r/go_final_project_git/internal/services"

	"github.com/go-chi/chi"
)

func Run() {
	defaultWebDir := "./internal/web/"
	defaultPort := "7540"
	dbDefaultPath := "internal/storage_db/scheduler.db"

	serverLink := services.SetServerLink(":", defaultPort)
	dbResultPath := services.CheckEnvDbVarOnExists(dbDefaultPath)

	log.Println("Checking DB on exists.")
	if err := services.CheckDBFileExists(dbResultPath); err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()

	fileServer := http.FileServer(http.Dir(defaultWebDir))
	log.Println("Creating handler.")
	r.Handle("/*", fileServer)

	log.Printf("Starting server on: '%s'\n", serverLink)

	srv := http.Server{
		Addr:         serverLink,
		Handler:      r,
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
	}

	log.Println("Listening...")
	if err := srv.ListenAndServe(); err != nil {
		log.Printf("Error starting server: %s", err.Error())
		return
	}
	log.Println("The server has stopped working.")
}
