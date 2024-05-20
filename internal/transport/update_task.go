package transport

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"fmt"

	"github.com/ra1nz0r/go_final_project/internal/config"
	"github.com/ra1nz0r/go_final_project/internal/database"
	"github.com/ra1nz0r/go_final_project/internal/services"
)

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	// Читаем данные из тела запроса.
	result, errBody := io.ReadAll(r.Body)
	if errBody != nil {
		config.LogErr.Error().Err(errBody).Msg("Cannot read from BODY.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Обрабатываем полученные данные из JSON и записываем в структуру.
	var task database.UpdateTaskParams
	if errUnm := json.Unmarshal(result, &task); errUnm != nil {
		services.ErrReturn(fmt.Errorf("can't deserialize: %w", errUnm), w)
		return
	}

	// Получаем путь из функции и подключаемся к датабазе.
	dbResPath, _ := services.CheckEnvDbVarOnExists(config.DbDefaultPath)
	db, errOpen := sql.Open("sqlite3", dbResPath)
	if errOpen != nil {
		config.LogErr.Fatal().Err(errOpen).Msg("Unable to connect to the database.")
	}

	// Проверяем корректность запроса для обновления параметров задачи в планировщике.
	if _, errFunc := services.NextDate(time.Now(), task.Date, task.Repeat); errFunc != nil {
		services.ErrReturn(fmt.Errorf("failed: %w", errFunc), w)
		return
	}

	// Если все данные введены корректно, то обновляем задачу в планировщике.
	queries := database.New(db)
	if errUpdate := queries.UpdateTask(context.Background(), task); errUpdate != nil {
		services.ErrReturn(fmt.Errorf("can't update task scheduler: %w", errUpdate), w)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(http.StatusAccepted)

	if _, errWrite := w.Write([]byte(`{}`)); errWrite != nil {
		config.LogErr.Error().Err(errWrite).Msg("Failed attempt WRITE response.")
		return
	}
}
