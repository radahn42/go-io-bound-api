package repo

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go-io-bound-api/internal/task"
	"sync"
	"testing"
	"time"
)

func TestTaskRepo(t *testing.T) {
	r := New()

	t.Run("Add and Get task", func(t *testing.T) {
		tk := &task.Task{
			ID:        "test-task",
			Status:    task.StatusPending,
			CreatedAt: time.Now(),
		}

		r.AddTask(tk)
		got, err := r.GetTask("test-task")

		require.NoError(t, err, "Should not return error")
		assert.Same(t, tk, got, "Should return same task instance")
	})

	t.Run("Get non-existent task", func(t *testing.T) {
		_, err := r.GetTask("non-existent")
		require.Error(t, err, "Should return error")
		assert.Contains(t, err.Error(), "not found", "Error should mention 'not found'")
	})

	t.Run("Delete existing task", func(t *testing.T) {
		tk := &task.Task{ID: "delete-me"}
		r.AddTask(tk)

		err := r.DeleteTask("delete-me")
		require.NoError(t, err, "Delete should be successful")

		_, err = r.GetTask("delete-me")
		assert.Error(t, err, "Task should be deleted")
	})

	t.Run("Delete non-existent task", func(t *testing.T) {
		err := r.DeleteTask("non-existent")
		require.Error(t, err, "Should return error")
		assert.Contains(t, err.Error(), "not found", "Error should mention 'not found'")
	})

	t.Run("Concurrent access", func(t *testing.T) {
		const numTasks = 100
		var wg sync.WaitGroup

		wg.Add(numTasks)
		for i := 0; i < numTasks; i++ {
			go func(i int) {
				defer wg.Done()
				tk := &task.Task{ID: fmt.Sprintf("task-%d", i)}
				r.AddTask(tk)
			}(i)
		}
		wg.Wait()

		for i := 0; i < numTasks; i++ {
			id := fmt.Sprintf("task-%d", i)
			_, err := r.GetTask(id)
			assert.NoError(t, err, "Task should exist")
		}
	})
}
