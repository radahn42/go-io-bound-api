package task

import (
	"github.com/google/uuid"
	"log/slog"
	"math/rand"
	"sync"
	"time"
)

// Status represents the current state of a task.
type Status string

var (
	StatusPending   Status = "pending"
	StatusRunning   Status = "running"
	StatusCompleted Status = "completed"
)

type SleepFuncType func(time.Duration)

var SleepFunc SleepFuncType = time.Sleep

// Task represents a singe long-running I/O bound operation.
type Task struct {
	ID          string
	Status      Status
	CreatedAt   time.Time
	StartedAt   time.Time
	CompletedAt time.Time

	mu sync.RWMutex
}

// New creates a new task with a unique ID and pending status.
func New() *Task {
	return &Task{
		ID:        uuid.NewString(),
		Status:    StatusPending,
		CreatedAt: time.Now(),
	}
}

func (t *Task) SimulateIOBoundWork() {
	const op = "task.SimulateIOBoundWork"
	log := slog.With(slog.String("op", op), slog.String("task_id", t.ID))

	log.Info("Starting simulation of I/O bound work")

	t.mu.Lock()
	t.Status = StatusRunning
	t.StartedAt = time.Now()
	t.mu.Unlock()

	minDuration := 3 * time.Minute
	maxDuration := 5 * time.Minute
	duration := minDuration + time.Duration(rand.Int63n(int64(maxDuration-minDuration+time.Nanosecond)))

	SleepFunc(duration)

	t.mu.Lock()
	t.Status = StatusCompleted
	t.CompletedAt = time.Now()
	t.mu.Unlock()

	slog.Info(
		"Completed simulation of I/O bound work",
		slog.String("duration_simulated", duration.String()),
		slog.Any("status", t.GetStatus()),
	)
}

// GetStatus returns the current status of the task.
func (t *Task) GetStatus() Status {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.Status
}

// GetCreatedAt returns the creation timestamp of the task.
func (t *Task) GetCreatedAt() time.Time {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.CreatedAt
}

// GetStartedAt returns the start timestamp of the task.
func (t *Task) GetStartedAt() time.Time {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.StartedAt
}

// GetCompletedAt returns the completion timestamp of the task.
func (t *Task) GetCompletedAt() time.Time {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.CompletedAt
}

// GetDuration calculates the current duration of the task if it's running or completed.
func (t *Task) GetDuration() string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	switch t.Status {
	case StatusPending:
		return ""
	case StatusRunning:
		return time.Since(t.StartedAt).Round(time.Second).String()
	case StatusCompleted:
		return t.CompletedAt.Sub(t.StartedAt).Round(time.Second).String()
	}
	return ""
}
