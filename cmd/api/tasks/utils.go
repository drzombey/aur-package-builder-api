package tasks

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func (at *ApiTask) logTaskError(taskId uuid.UUID, taskName string, message string) {
	log.Errorf("[Task: %s][ID: %s] %s", taskName, taskId.String(), message)
}

func (at *ApiTask) logTaskCompleted(taskId uuid.UUID, taskName string) {
	log.Infof("[Task: %s][ID: %s] Completed successfully", taskName, taskId.String())
}
