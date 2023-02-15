package pkg

import (
	"context"
	"testing"
	"time"
	// . "github.com/onsi/ginkgo"
	// . "github.com/onsi/gomega"
)

func TestCheckIsntanceRunning(t *testing.T) {
	t.Logf("test again instance running state")
}

func TestCheckIsntanceStopped(t *testing.T) {
	t.Logf("test again instance stopped state")
}

func TestCheckExpectedTime(t *testing.T) {
	t.Logf("test trigger time against working time")

	expectedOK := true
	actualTime := time.Now().Add(8 * time.Hour)

	working, err := checkExpectedTime(context.TODO(), actualTime)
	if err != nil {
		t.Errorf("format time failed: %v", err)
	}

	if working != expectedOK {
		t.Errorf("wrong time check current time is %s, while shows not in work", actualTime.String())
	}
}

func TestCheckNotExpectedTime(t *testing.T) {
	t.Logf("test trigger time against non-working time")

	expectedOK := false
	actualTime := time.Now().Add(19 * time.Hour)

	working, err := checkExpectedTime(context.TODO(), actualTime)
	if err != nil {
		t.Errorf("format time failed: %v", err)
	}

	if working != expectedOK {
		t.Errorf("wrong time check current time is %s, while shows in work", actualTime.String())
	}
}
