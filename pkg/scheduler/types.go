package scheduler

type TasksScheduler struct {
	taskChannels []chan struct{}
}

func NewTasksScheduler() *TasksScheduler {
	return &TasksScheduler{}
}
