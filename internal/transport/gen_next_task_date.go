package transport

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/ra1nz0r/go_final_project/internal/config"
	"github.com/ra1nz0r/go_final_project/internal/database"
	"github.com/ra1nz0r/go_final_project/internal/services"
)

func GeneratedNextDate(w http.ResponseWriter, r *http.Request) {
	// Получаем путь из функции и подключаемся к базе данных.
	dbResPath, _ := services.CheckEnvDbVarOnExists(config.DbDefaultPath)
	db, errOpen := sql.Open("sqlite3", dbResPath)
	if errOpen != nil {
		config.LogErr.Fatal().Err(errOpen).Msg("Unable to connect to the database.")
	}

	// Получаем задачу по ID и возвращаем ошибку, если её нет в базе данных.
	queries := database.New(db)
	taskGeted, errGeted := queries.GetTask(context.Background(), r.URL.Query().Get("id"))
	if errGeted != nil {
		services.ErrReturn(make(map[string]string), fmt.Sprintf("The ID you entered does not exist: %v", errGeted), w)
		return
	}

	switch {
	case taskGeted.Repeat == "": // Одноразовая задача с пустым полем REPEAT удаляется.
		if errDel := queries.DeleteTask(context.Background(), taskGeted.ID); errDel != nil {
			services.ErrReturn(make(map[string]string), fmt.Sprintf("Failed delete: %v", errDel), w)
			return
		}
	default: // В остальных случаях, расчитывается и записывается новая дата для задачи вместо старой.
		newDate, errFunc := services.NextDate(time.Now(), taskGeted.Date, taskGeted.Repeat)
		if errFunc != nil {
			services.ErrReturn(make(map[string]string), fmt.Sprintf("Failed: %v", errFunc), w)
			return
		}

		var task database.UpdateDateTaskParams
		task.ID = taskGeted.ID
		task.Date = newDate
		if errUpd := queries.UpdateDateTask(context.Background(), task); errUpd != nil {
			services.ErrReturn(make(map[string]string), fmt.Sprintf("Failed update task: %v", errUpd), w)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(http.StatusOK)

	if _, errWrite := w.Write([]byte(`{}`)); errWrite != nil {
		config.LogErr.Error().Err(errWrite).Msg("Failed attempt WRITE response.")
		return
	}

}
