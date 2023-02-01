package scheduler

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestAppendTaskChannelToList(t *testing.T) {
	ts := &TasksScheduler{}

	taskChannel1 := make(chan struct{})
	taskChannel2 := make(chan struct{})

	ts.appendTaskChannelToList(&taskChannel1)
	ts.appendTaskChannelToList(&taskChannel2)

	if len(ts.taskChannels) != 2 {
		t.Errorf("Expected taskChannels to have length 2, but got %d", len(ts.taskChannels))
	}

	if ts.taskChannels[0] != taskChannel1 {
		t.Errorf("Expected taskChannels[0] to be %v, but got %v", taskChannel1, ts.taskChannels[0])
	}

	if ts.taskChannels[1] != taskChannel2 {
		t.Errorf("Expected taskChannels[1] to be %v, but got %v", taskChannel2, ts.taskChannels[1])
	}
}

func TestLogTaskStartingInfo(t *testing.T) {
	ts := &TasksScheduler{}

	var buf bytes.Buffer
	logrus.SetOutput(&buf)
	defer logrus.SetOutput(os.Stderr)

	taskName := "test-task"
	ts.logTaskStartingInfo(taskName)

	expectedLog := fmt.Sprintf(`time="\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}.*" level=info msg="\[%s\] is starting"`, taskName)

	match, _ := regexp.MatchString(expectedLog, buf.String())
	if !match {
		t.Errorf("Expected log message to match %q, but got %q", expectedLog, buf.String())
	}
}

func TestLogTaskInitInfo(t *testing.T) {
	ts := &TasksScheduler{}

	var buf bytes.Buffer
	logrus.SetOutput(&buf)
	defer logrus.SetOutput(os.Stderr)

	taskName := "test-task"
	ts.logTaskInitInfo(taskName)

	expectedLog := fmt.Sprintf(`time="\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}.*" level=info msg="\[%s\] is initialized"`, taskName)

	match, _ := regexp.MatchString(expectedLog, buf.String())
	if !match {
		t.Errorf("Expected log message to match %q, but got %q", expectedLog, buf.String())
	}
}

func TestLogTaskStoppingInfo(t *testing.T) {
	ts := &TasksScheduler{}

	var buf bytes.Buffer
	logrus.SetOutput(&buf)
	defer logrus.SetOutput(os.Stderr)

	taskName := "test-task"
	ts.logTaskStoppingInfo(taskName)

	expectedLog := fmt.Sprintf(`time="\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}.*" level=info msg="\[%s\] is stopping"`, taskName)

	match, _ := regexp.MatchString(expectedLog, buf.String())
	if !match {
		t.Errorf("Expected log message to match %q, but got %q", expectedLog, buf.String())
	}
}
