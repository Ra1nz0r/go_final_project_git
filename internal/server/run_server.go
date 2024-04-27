package server

import (
	"log"
	"net/http"

	"github.com/ra1nz0r/go_final_project_git/internal/services"

	"github.com/go-chi/chi"
)

func Run() {
	r := chi.NewRouter()

	defaultWebDir := "./internal/web/"
	defaultPort := "7540"
	dbDefaultPath := "internal/storage_db/scheduler.db"

	serverLink := services.SetServerLink(":", defaultPort)
	dbResultPath := services.CheckEnvDbVarOnExists(dbDefaultPath)

	log.Println("Checking DB on exists.")
	if err := services.CheckDBFileExists(dbResultPath); err != nil {
		log.Fatal(err)
	}

	fileServer := http.FileServer(http.Dir(defaultWebDir))
	log.Println("Creating handler.")
	r.Handle("/*", fileServer)

	log.Printf("Starting server on: '%s'\n", serverLink)
	log.Println("Listening...")
	if err := http.ListenAndServe(serverLink, r); err != nil {
		log.Printf("Error starting server: %s", err.Error())
		return
	}
}
