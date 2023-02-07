package scheduler

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestScheduleTask(t *testing.T) {
	ts := &TasksScheduler{}
	taskCalled := false

	taskFn := func(taskId uuid.UUID, taskName string) {
		taskCalled = true
	}

	ts.ScheduleTask(taskFn, 1, "Test Task")

	time.Sleep(2 * time.Second)

	if !taskCalled {
		t.Errorf("Task was not called")
	}
}

func TestStopTask(t *testing.T) {
	ts := &TasksScheduler{}
	taskCalled := false

	taskFn := func(taskId uuid.UUID, taskName string) {
		taskCalled = true
	}

	ts.ScheduleTask(taskFn, 1, "Test Task")
	ts.ShutdownAllTasks()

	time.Sleep(2 * time.Second)

	if taskCalled {
		t.Errorf("Task was called even after stopping it")
	}
}
