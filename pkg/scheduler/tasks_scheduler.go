package scheduler

import (
	"time"

	"github.com/google/uuid"
)

func (ts *TasksScheduler) ScheduleTask(taskFn func(taskId uuid.UUID, taskName string), timeInSeconds int, taskName string) {
	ticker := time.NewTicker(time.Duration(timeInSeconds) * time.Second)
	quit := make(chan struct{})
	taskId := uuid.New()
	go func() {
		ts.logTaskInitInfo(taskName)
		for {
			select {
			case <-ticker.C:
				ts.logTaskStartingInfo(taskName)
				taskFn(taskId, taskName)
			case <-quit:
				ts.logTaskStoppingInfo(taskName)
				ticker.Stop()
				return
			}
		}
	}()

	ts.appendTaskChannelToList(&quit)
}

func (ts *TasksScheduler) ShutdownAllTasks() {
	for _, task := range ts.taskChannels {
		close(task)
	}
}
