package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ra1nz0r/go_final_project/internal/config"
	"github.com/ra1nz0r/go_final_project/internal/services"
	"github.com/ra1nz0r/go_final_project/internal/transport"

	"github.com/go-chi/chi"
)

func Run() {
	serverLink, boolValue := services.SetServerLink(":", config.DefaultPort)
	if boolValue {
		config.LogInfo.Info().Msg("'TODO_PORT' exitst in '.env' file. Changing default PORT.")
	}
	dbResultPath, boolValue := services.CheckEnvDbVarOnExists(config.DbDefaultPath)
	if boolValue {
		config.LogInfo.Info().Msg("'TODO_DBFILE' exitst in '.env' file. Changing default PATH.")
	}

	config.LogInfo.Info().Msg("Checking DB on exists.")
	if errCheck := services.CheckDBFileExists(dbResultPath); errCheck != nil {
		config.LogErr.Fatal().Err(errCheck).Msgf("Cannot check DB on exists.")
	}

	r := chi.NewRouter()

	fileServer := http.FileServer(http.Dir(config.DefaultWebDir))
	config.LogInfo.Info().Msg("Running handlers.")
	r.Handle("/*", fileServer)

	r.Get("/api/nextdate", transport.NextDateHand)

	r.Get("/api/tasks", transport.CheckAuth(transport.UpcomingTasksWithSearch))

	r.Post("/api/task/done", transport.CheckAuth(transport.GeneratedNextDate))

	r.Post("/api/signin", transport.LoginAuth)

	r.Delete("/api/task", transport.CheckAuth(transport.DeleteTaskScheduler))
	r.Get("/api/task", transport.CheckAuth(transport.GetTaskByID))
	r.Post("/api/task", transport.CheckAuth(transport.AddSchedulerTask))
	r.Put("/api/task", transport.CheckAuth(transport.UpdateTask))

	config.LogInfo.Info().Msgf("Starting server on: '%s'", serverLink)

	srv := http.Server{
		Addr:         serverLink,
		Handler:      r,
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
	}

	go func() {
		if errListn := srv.ListenAndServe(); !errors.Is(errListn, http.ErrServerClosed) {
			config.LogErr.Fatal().Err(errListn).Msg("HTTP server error.")
		}
		config.LogInfo.Info().Msg("Stopped serving new connections.")
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if errShut := srv.Shutdown(shutdownCtx); errShut != nil {
		config.LogErr.Fatal().Err(errShut).Msg("HTTP shutdown error")
	}
	config.LogInfo.Info().Msg("Graceful shutdown complete.")
}
