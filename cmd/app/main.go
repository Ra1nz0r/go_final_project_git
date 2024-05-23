package main

import (
	_ "github.com/joho/godotenv/autoload"
	_ "modernc.org/sqlite"

	"github.com/ra1nz0r/go_final_project/internal/server"
)

func main() {
	server.Run()
}
