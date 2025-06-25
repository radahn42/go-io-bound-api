package task

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestNewTaskDefaults(t *testing.T) {
	task := New()

	assert.NotEmpty(t, task.ID, "ID should not be empty")
	assert.Equal(t, StatusPending, task.GetStatus(), "Default status should be pending")
	assert.False(t, task.GetCreatedAt().IsZero(), "CreatedAt should be set")
}

func TestTaskStatusTransitions(t *testing.T) {
	task := New()

	task.mu.Lock()
	task.Status = StatusRunning
	task.StartedAt = time.Now().Add(-time.Second)
	task.mu.Unlock()

	assert.Equal(t, StatusRunning, task.GetStatus(), "Status should be running")
	assert.False(t, task.GetStartedAt().IsZero(), "StartedAt should be set")

	task.mu.Lock()
	task.Status = StatusCompleted
	task.CompletedAt = time.Now()
	task.mu.Unlock()

	assert.Equal(t, StatusCompleted, task.GetStatus(), "Status should be completed")
	assert.False(t, task.GetCompletedAt().IsZero(), "CompletedAt should be set")
}

func TestTaskGetDuration(t *testing.T) {
	task := New()

	assert.Empty(t, task.GetDuration(), "Duration should be empty for pending task")

	task.mu.Lock()
	task.Status = StatusRunning
	task.StartedAt = time.Now().Add(-3 * time.Second)
	task.mu.Unlock()

	assert.NotEmpty(t, task.GetDuration(), "Duration should not be empty for running task")
	assert.Contains(t, task.GetDuration(), "3s", "Duration should reflect running time")

	task.mu.Lock()
	task.Status = StatusCompleted
	task.CompletedAt = task.StartedAt.Add(5 * time.Second)
	task.mu.Unlock()

	assert.Equal(t, "5s", task.GetDuration(), "Duration should be 5s for completed task")
}

func TestSimulateIOBoundWork(t *testing.T) {
	origSleep := SleepFunc
	defer func() {
		SleepFunc = origSleep
	}()

	var sleepCalled bool
	var sleepDuration time.Duration
	SleepFunc = func(d time.Duration) {
		sleepCalled = true
		sleepDuration = d
	}

	task := New()
	task.SimulateIOBoundWork()

	assert.True(t, sleepCalled, "Sleep should be called")
	assert.GreaterOrEqual(t, sleepDuration, 3*time.Minute, "Sleep should be at least 3 minutes")
	assert.LessOrEqual(t, sleepDuration, 5*time.Minute, "Sleep should be at most 5 minutes")
	assert.Equal(t, StatusCompleted, task.GetStatus(), "Status should be completed")
	assert.False(t, task.GetStartedAt().IsZero(), "StartedAt should be set")
	assert.False(t, task.GetCompletedAt().IsZero(), "CompletedAt should be set")
}

func TestConcurrentTaskExecution(t *testing.T) {
	origSleep := SleepFunc
	defer func() {
		SleepFunc = origSleep
	}()

	var sleepCalls int
	var mu sync.Mutex
	SleepFunc = func(d time.Duration) {
		mu.Lock()
		sleepCalls++
		mu.Unlock()
	}

	var wg sync.WaitGroup
	wg.Add(10)
	tasks := make([]*Task, 10)
	for i := 0; i < 10; i++ {
		tasks[i] = New()
		go func(tsk *Task) {
			defer wg.Done()
			tsk.SimulateIOBoundWork()
		}(tasks[i])
	}
	wg.Wait()

	assert.Equal(t, 10, sleepCalls, "Should have 10 sleep calls")
	for _, task := range tasks {
		assert.Equal(t, StatusCompleted, task.GetStatus(), "All tasks should be completed")
	}
}
