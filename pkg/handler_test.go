package pkg

import "testing"

func TestCheckIsntanceRunning(t *testing.T) {
	t.Logf("test again instance running state")
}

func TestCheckIsntanceStopped(t *testing.T) {
	t.Logf("test again instance stopped state")
}

func TestCheckExpectedTime(t *testing.T) {
	t.Logf("test trigger time against working time")
}

func TestCheckNotExpectedTime(t *testing.T) {
	t.Logf("test trigger time against non-working time")
}
