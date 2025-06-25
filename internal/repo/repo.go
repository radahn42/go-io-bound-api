package repo

import (
	"fmt"
	"go-io-bound-api/internal/task"
	"sync"
)

// TaskRepo manages the in-memory storage of tasks.
type TaskRepo struct {
	mu    sync.RWMutex
	tasks map[string]*task.Task
}

// New creates a new in-memory task repository
func New() *TaskRepo {
	return &TaskRepo{
		tasks: make(map[string]*task.Task),
	}
}

// AddTask adds a new task to the repository
func (r *TaskRepo) AddTask(t *task.Task) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tasks[t.ID] = t
}

// GetTask retrieves a task by its ID.
func (r *TaskRepo) GetTask(id string) (*task.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.tasks[id]
	if !ok {
		return nil, fmt.Errorf("task with ID %s not found", id)
	}
	return t, nil
}

// DeleteTask removes a task from the repository.
func (r *TaskRepo) DeleteTask(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.tasks[id]; !ok {
		return fmt.Errorf("task with ID %s not found", id)
	}
	delete(r.tasks, id)
	return nil
}
