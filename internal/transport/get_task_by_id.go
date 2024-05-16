package transport

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"fmt"

	"github.com/ra1nz0r/go_final_project/internal/config"
	"github.com/ra1nz0r/go_final_project/internal/database"
	"github.com/ra1nz0r/go_final_project/internal/services"
)

func GetTaskById(w http.ResponseWriter, r *http.Request) {
	// Получаем путь из функции и подключаемся к датабазе.
	dbResPath, _ := services.CheckEnvDbVarOnExists(config.DbDefaultPath)
	db, errOpen := sql.Open("sqlite3", dbResPath)
	if errOpen != nil {
		config.LogErr.Fatal().Err(errOpen).Msg("Unable to connect to the database.")
	}

	// Получаем задачу из планировщика при GET запросе в виде "/api/task?id=185".
	queries := database.New(db)
	idGeted, errGeted := queries.GetTask(context.Background(), r.URL.Query().Get("id"))
	if errGeted != nil {
		services.ErrReturn(make(map[string]string), fmt.Sprintf("The ID you entered does not exist: %v", errGeted), w)
		return
	}

	// Оборачиваем полученные данные в JSON и готовим к выводу,
	// ответ в виде: {"id": "айди","date": "дата","title": "заголовок","comment": "коммент","repeat": "условия повторения"}.
	jsonResp, errJson := json.Marshal(idGeted)
	if errJson != nil {
		config.LogErr.Error().Err(errJson).Msg("Failed attempt json-marshal response.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(http.StatusOK)

	if _, errWrite := w.Write(jsonResp); errWrite != nil {
		config.LogErr.Error().Err(errWrite).Msg("Failed attempt WRITE response.")
		return
	}

}
