package docker

import "testing"

func TestEnsureImage(t *testing.T) {
	c, err := NewContainerController()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	err = c.EnsureImage("alpine")

	if err != nil {
		t.Error(err)
	}
}
