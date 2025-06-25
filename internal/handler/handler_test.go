package handler

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go-io-bound-api/internal/repo"
	"go-io-bound-api/internal/task"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCreateTask_MockSimulateIOBoundWork(t *testing.T) {
	origSleep := task.SleepFunc
	defer func() {
		task.SleepFunc = origSleep
	}()

	repository := repo.New()
	handler := New(repository)
	r := chi.NewRouter()
	r.Post("/tasks", handler.CreateTask)
	r.Get("/tasks/{id}", handler.GetTaskStatus)
	r.Delete("/tasks/{id}", handler.DeleteTask)

	t.Run("Create task", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/tasks", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusAccepted, w.Code, "Status should be 202 Accepted")

		var resp map[string]any
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))

		id, ok := resp["id"].(string)
		require.True(t, ok, "Response should contain id")
		assert.NotEmpty(t, id, "Task ID should not be empty")
	})

	t.Run("Get task status", func(t *testing.T) {
		tk := &task.Task{
			ID:        "test-task",
			Status:    task.StatusRunning,
			CreatedAt: time.Now(),
			StartedAt: time.Now().Add(-time.Minute),
		}
		repository.AddTask(tk)

		req := httptest.NewRequest(http.MethodGet, "/tasks/test-task", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Status should be 200 OK")

		var resp map[string]any
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))

		assert.Equal(t, "test-task", resp["id"], "Task ID should match")
		assert.Equal(t, "running", resp["status"], "Status should be running")
		assert.NotEmpty(t, resp["duration"], "Duration should be present")
	})

	t.Run("Get non-existent task", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/tasks/non-existent", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code, "Status should be 404 Not Found")
	})

	t.Run("Delete task", func(t *testing.T) {
		tk := &task.Task{ID: "delete-me"}
		repository.AddTask(tk)

		req := httptest.NewRequest(http.MethodDelete, "/tasks/delete-me", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code, "Status should be 204 No Content")

		_, err := repository.GetTask("delete-me")
		assert.Error(t, err, "Task should be deleted")
	})

	t.Run("Delete non-existent task", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/tasks/non-existent", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code, "Status should be 404 Not Found")
	})

	t.Run("Full workflow", func(t *testing.T) {
		createReq := httptest.NewRequest(http.MethodPost, "/tasks", nil)
		createRes := httptest.NewRecorder()
		r.ServeHTTP(createRes, createReq)
		require.Equal(t, http.StatusAccepted, createRes.Code)

		var createResp map[string]string
		require.NoError(t, json.Unmarshal(createRes.Body.Bytes(), &createResp))
		taskID := createResp["id"]
		require.NotEmpty(t, taskID)

		time.Sleep(10 * time.Millisecond)

		statusReq := httptest.NewRequest(http.MethodGet, "/tasks/"+taskID, nil)
		statusRes := httptest.NewRecorder()
		r.ServeHTTP(statusRes, statusReq)
		require.Equal(t, http.StatusOK, statusRes.Code)

		var statusResp map[string]any
		require.NoError(t, json.Unmarshal(statusRes.Body.Bytes(), &statusResp))
		assert.Contains(t, []string{"completed", "running"}, statusResp["status"], "Status should be completed or running")

		deleteReq := httptest.NewRequest(http.MethodDelete, "/tasks/"+taskID, nil)
		deleteRes := httptest.NewRecorder()
		r.ServeHTTP(deleteRes, deleteReq)
		assert.Equal(t, http.StatusNoContent, deleteRes.Code)

		_, err := repository.GetTask(taskID)
		assert.Error(t, err, "Task should be deleted")
	})
}
