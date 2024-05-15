package services

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ra1nz0r/go_final_project/internal/config"
	"github.com/ra1nz0r/go_final_project/internal/models"
)

// Создает адрес запуска сервера и изменяет порт прослушивания по умолчанию,
// если переменная «TODO _PORT» существует в «.env».
// Переменная bool используется один раз, для вывода сообщения о существовании перменной «TODO _PORT» в «.env»
// и изменении стандартного порта при запуске сервера, в остальных случаях пропускается.
func SetServerLink(address string, port string) (string, bool) {
	bool := false
	if portFromEnv, exists := os.LookupEnv("TODO_PORT"); exists && portFromEnv != "" {
		port = portFromEnv
		bool = true
	}
	return address + port, bool
}

// Изменяет путь по умолчанию к базе данных на «TODO_DBFILE», если переменная существует в «.env».
// Переменная bool используется один раз, для вывода сообщения о существовании перменной «TODO_DBFILE» в «.env»
// и изменении стандартного пути датабазы при запуске сервера, в остальных случаях пропускается.
func CheckEnvDbVarOnExists(dbDefaultPath string) (string, bool) {
	bool := false
	if dbPathFromEnv := os.Getenv("TODO_DBFILE"); dbPathFromEnv != "" {
		dbDefaultPath, bool = dbPathFromEnv, true
	}
	return dbDefaultPath, bool
}

// Проверка существования DB.
// Создание папки для хранения DB, файла «.db» и TABLE.
func CheckDBFileExists(resPath string) error {
	if _, errStat := os.Stat(resPath); errStat != nil {
		if os.IsNotExist(errStat) {

			// Создание папки хранения для базы данных.
			folderDb := filepath.Dir(resPath)
			if errMkDir := os.Mkdir(folderDb, 0777); errMkDir != nil {
				return fmt.Errorf("failed: cannot create folder: %v", errMkDir)
			}

			config.LogInfo.Info().Msgf("Creating %s and TABLE.", filepath.Base(resPath))
			ctx := context.Background()
			db, errOpen := sql.Open("sqlite3", resPath)
			if errOpen != nil {
				return fmt.Errorf("failed: cannot open db: %v", errOpen)
			}

			// Создание TABLE.
			if _, errCreate := db.ExecContext(ctx, models.Ddl); errCreate != nil {
				return fmt.Errorf("failed: cannot create table db: %v", errCreate)
			}
			return nil
		}
	}
	config.LogInfo.Info().Msgf("Database %s exists.", filepath.Base(resPath))
	return nil
}
