package handler

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"go-io-bound-api/internal/dto"
	"go-io-bound-api/internal/repo"
	"go-io-bound-api/internal/task"
	"log/slog"
	"net/http"
)

// TaskHandler provides HTTP handlers for task operations.
type TaskHandler struct {
	repo *repo.TaskRepo
}

// New creates a new TaskHandler
func New(repo *repo.TaskRepo) *TaskHandler {
	return &TaskHandler{
		repo: repo,
	}
}

// CreateTask handles the creation of a new task.
//
// POST /tasks
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	const op = "handler.CreateTask"
	log := slog.With(slog.String("op", op))

	newTask := task.New()
	h.repo.AddTask(newTask)

	log.Info("New task created and accepted", slog.String("task_id", newTask.ID))

	go func() {
		newTask.SimulateIOBoundWork()
	}()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"id":      newTask.ID,
		"status":  string(newTask.GetStatus()),
		"message": "Task accepted and processing",
	})
}

// GetTaskStatus retrieves the status of a task.
//
// GET /tasks/{id}
func (h *TaskHandler) GetTaskStatus(w http.ResponseWriter, r *http.Request) {
	const op = "handler.GetTaskStatus"
	log := slog.With(slog.String("op", op))

	id := chi.URLParam(r, "id")
	if id == "" {
		log.Warn("Task ID is missing in URL for GetTaskStatus")
		http.Error(w, "Task ID is required", http.StatusBadRequest)
		return
	}

	t, err := h.repo.GetTask(id)
	if err != nil {
		log.Warn(
			"Attempted to get non-existent task",
			slog.String("task_id", id),
			slog.Any("error", err),
		)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	log.Info(
		"Retrieving task status",
		slog.String("task_id", id),
		slog.Any("status", t.GetStatus()),
	)

	resp := dto.TaskResponse{
		ID:          t.ID,
		Status:      t.GetStatus(),
		CreatedAt:   t.GetCreatedAt(),
		StartedAt:   t.GetStartedAt(),
		CompletedAt: t.GetCompletedAt(),
		Duration:    t.GetDuration(),
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// DeleteTask handles the deletion of a task.
//
// DELETE /tasks/{id}
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	const op = "handler.DeleteTask"
	log := slog.With(slog.String("op", op))

	id := chi.URLParam(r, "id")
	if id == "" {
		log.Warn("Task ID is missing in URL for DeleteTask")
		http.Error(w, "Task ID is required", http.StatusBadRequest)
		return
	}

	err := h.repo.DeleteTask(id)
	if err != nil {
		log.Warn(
			"Attempted to delete non-existent task",
			slog.String("task_id", id),
			slog.Any("error", err),
		)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	log.Info("Task deleted successfully", slog.String("task_id", id))
	w.WriteHeader(http.StatusNoContent)
}
