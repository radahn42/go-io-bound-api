package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go-io-bound-api/internal/handler"
	"go-io-bound-api/internal/repo"
	"log"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	taskRepo := repo.New()
	taskHandler := handler.New(taskRepo)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/tasks", taskHandler.CreateTask)
	r.Get("/tasks/{id}", taskHandler.GetTaskStatus)
	r.Delete("/tasks/{id}", taskHandler.DeleteTask)

	port := ":8080"
	logger.Info("Server starting", slog.String("port", port))
	log.Fatal(http.ListenAndServe(port, r))
}
