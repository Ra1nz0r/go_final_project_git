package main

import (
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"

	"github.com/ra1nz0r/go_final_project_git/internal/server"
)

func main() {

	server.Run()
}
