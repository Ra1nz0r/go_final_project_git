package services

import (
	"context"
	"database/sql"
	"log"
	"os"
	"path/filepath"

	"github.com/ra1nz0r/go_final_project_git/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

// Creates a server startup address and changes the default
// listening port, if the "TODO _PORT" variable exists in '.env'.
func SetServerLink(address string, port string) string {
	if TodoPort, exists := os.LookupEnv("TODO_PORT"); exists && TodoPort != "" {
		log.Println("'TODO_PORT' exitst in '.env' file. Changing default PORT.")
		port = TodoPort
	}
	return address + port
}

// Changes the default path to the database on "TODO_DBFILE", if a variable exists in '.env'.
func CheckEnvDbVarOnExists(dbDefaultPath string) string {
	if dbPathFromSetting := os.Getenv("TODO_DBFILE"); dbPathFromSetting != "" {
		log.Println("'TODO_DBFILE' exitst in '.env' file. Changing default PATH.")
		dbDefaultPath = dbPathFromSetting
	}
	return dbDefaultPath
}

// Checking for the existence of a DB.
// Creating a folder to store DB, '.db' file and TABLE.
func CheckDBFileExists(resPath string) error {
	if _, err := os.Stat(resPath); err != nil {
		if os.IsNotExist(err) {

			// Creating a storage folder for the database.
			folderDb := filepath.Dir(resPath)
			if err := os.Mkdir(folderDb, 0777); err != nil {
				log.Println(err)
			}

			log.Printf("Creating %s and TABLE.", filepath.Base(resPath))
			ctx := context.Background()
			db, err := sql.Open("sqlite3", resPath)
			if err != nil {
				return err
			}

			// Creating TABLE.
			if _, err := db.ExecContext(ctx, models.Ddl); err != nil {
				return err
			}
			return nil
		}
	}
	log.Println("Database 'scheduler.db' exists.")
	return nil
}
