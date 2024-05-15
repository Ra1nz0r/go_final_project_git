package transport

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/ra1nz0r/go_final_project/internal/config"
	"github.com/ra1nz0r/go_final_project/internal/database"
	"github.com/ra1nz0r/go_final_project/internal/services"
)

func DeleteTaskScheduler(w http.ResponseWriter, r *http.Request) {
	// Получаем путь из функции и подключаемся к датабазе.
	dbResPath, _ := services.CheckEnvDbVarOnExists(config.DbDefaultPath)
	db, errOpen := sql.Open("sqlite3", dbResPath)
	if errOpen != nil {
		config.LogErr.Fatal().Err(errOpen).Msg("Unable to connect to the database.")
	}

	// Проверям существование задачи и возвращаем ошибку, если её нет в базе данных.
	queries := database.New(db)
	_, errGeted := queries.GetTask(context.Background(), r.URL.Query().Get("id"))
	if errGeted != nil {
		services.ErrReturn(make(map[string]string), fmt.Sprintf("The ID you entered does not exist: %v", errGeted), w)
		return
	}

	// Удаляем задачу из базы данных, при DELETE запросе в виде "/api/task?id=185".
	if errDel := queries.DeleteTask(context.Background(), r.URL.Query().Get("id")); errDel != nil {
		services.ErrReturn(make(map[string]string), fmt.Sprintf("Failed delete: %v", errDel), w)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(http.StatusOK)

	if _, errWrite := w.Write([]byte(`{}`)); errWrite != nil {
		config.LogErr.Error().Err(errWrite).Msg("Failed attempt WRITE response.")
		return
	}
}
