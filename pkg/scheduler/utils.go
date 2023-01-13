package scheduler

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

func (ts *TasksScheduler) appendTaskChannelToList(taskChannel *chan struct{}) {
	ts.taskChannels = append(ts.taskChannels, *taskChannel)
}

func (ts *TasksScheduler) logTaskStartingInfo(taskName string) {
	logrus.Info(fmt.Sprintf("[%s] is starting", taskName))
}

func (ts *TasksScheduler) logTaskInitInfo(taskName string) {
	logrus.Info(fmt.Sprintf("[%s] is initialized", taskName))
}

func (ts *TasksScheduler) logTaskStoppingInfo(taskName string) {
	logrus.Info(fmt.Sprintf("[%s] is stopping", taskName))
}
