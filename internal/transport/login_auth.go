package transport

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"fmt"

	"github.com/ra1nz0r/go_final_project/internal/config"
	"github.com/ra1nz0r/go_final_project/internal/services"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Password string `json:"password"`
}

func LoginAuth(w http.ResponseWriter, r *http.Request) {
	// Читаем данные из тела запроса.
	result, errBody := io.ReadAll(r.Body)
	if errBody != nil {
		config.LogErr.Error().Err(errBody).Msg("Cannot read from BODY.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Обрабатываем полученные данные из JSON и записываем в структуру.
	var u User
	if errUnm := json.Unmarshal(result, &u); errUnm != nil {
		services.ErrReturn(fmt.Errorf("can't deserialize: %w", errUnm), w)
		return
	}

	// Проверяем существование переменной "TODO_PASSWORD" в ".env".
	// В случае успеха записываем в результат хэш, в противном ошибку.
	passFromEnv := os.Getenv("TODO_PASSWORD")
	respResult := make(map[string]string)
	switch {
	case passFromEnv == u.Password:
		passHash, errCrypt := bcrypt.GenerateFromPassword([]byte(passFromEnv), bcrypt.DefaultCost)
		if errCrypt != nil {
			services.ErrReturn(fmt.Errorf("failed to generate password hash: %w", errCrypt), w)
		}
		respResult["token"] = string(passHash)
	default:
		respResult["error"] = "Incorrect password."
	}

	// Оборачиваем полученные данные в JSON и готовим к выводу,
	// ответ в виде: {"token/error":"hash/txt_error"}.
	jsonResp, errJSON := json.Marshal(respResult)
	if errJSON != nil {
		config.LogErr.Error().Err(errJSON).Msg("Failed attempt json-marshal response.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(http.StatusAccepted)

	if _, errWrite := w.Write(jsonResp); errWrite != nil {
		config.LogErr.Error().Err(errWrite).Msg("Failed attempt WRITE response.")
		return
	}
}
