package docker

import "testing"

func TestContainerRun(t *testing.T) {
	c, err := NewContainerController()

	if err != nil {
		t.Error(err)
	}

	containerId, err := c.RunContainer("alpine", []string{"echo", "hello world"}, []VolumeMount{})

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if containerId == "" {
		t.Error("Expected containerId is filled; received none")
	}

	c.ContainerById(containerId)
}
