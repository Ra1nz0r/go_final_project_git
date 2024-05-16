package services

import (
	"encoding/json"
	"net/http"

	"github.com/ra1nz0r/go_final_project/internal/config"
)

// Оборачивает ошибки в JSON и возвращает в формате {"error":"ваш текст для ошибки"}.
func ErrReturn(result map[string]string, err string, w http.ResponseWriter) {
	result["error"] = err
	jsonResp, errJson := json.Marshal(result)
	if errJson != nil {
		config.LogErr.Error().Err(errJson).Msg("Failed attempt json-marshal response.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(http.StatusBadRequest)

	if _, errWrite := w.Write(jsonResp); errWrite != nil {
		config.LogErr.Error().Err(errWrite).Msgf("Failed attempt WRITE response.")
		return
	}
}
