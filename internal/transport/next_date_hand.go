package transport

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ra1nz0r/go_final_project/internal/config"
	"github.com/ra1nz0r/go_final_project/internal/services"
)

// Обработчик для GET запросов и вывода следующей даты в текстовом формате
// для задачи в планировщике, после обработки функциeй NextDate.
func NextDateHand(w http.ResponseWriter, r *http.Request) {
	// Обрабатываем введенную дату.
	todayDate, errPars := time.Parse("20060102", r.URL.Query().Get("now"))
	if errPars != nil {
		errorMsg := fmt.Sprintf("Failed: incorrect DATE: %v", errPars)
		http.Error(w, errorMsg, http.StatusBadRequest)
		return
	}

	// Вычисление следующей даты, подробнее в описании NextDate.
	res, errFunc := services.NextDate(todayDate, r.URL.Query().Get("date"), r.URL.Query().Get("repeat"))
	if errFunc != nil {
		http.Error(w, errFunc.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")

	w.WriteHeader(http.StatusOK)

	if _, errWrite := w.Write([]byte(res)); errWrite != nil {
		config.LogErr.Error().Err(errWrite).Msgf("Failed attempt WRITE response.")
		return
	}
}
