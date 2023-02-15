package pkg

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	// . "github.com/onsi/ginkgo"
	// . "github.com/onsi/gomega"
)

func TestCheckInstanceRunning(t *testing.T) {
	t.Logf("test again instance running state")
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		t.Errorf("fail to load the config: %v", err)
	}

	// pls enter your own instance id and profile
	client := ec2.NewFromConfig(cfg)
	fakeInstanceID := "i-04b465ab13755d687"
	if err := stopInstance(context.TODO(), fakeInstanceID, client, true); err != nil {
		if !strings.Contains(err.Error(), "api error DryRunOperation: Request would have succeeded, but DryRun flag is set") {
			t.Errorf("fail to stop the instance: %v", err)
		}
		t.Skip()
	}

}

func TestCheckIsntanceStopped(t *testing.T) {
	t.Logf("test again instance stopped state")
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		t.Errorf("fail to load the config: %v", err)
	}

	// pls enter your own instance id and profile
	client := ec2.NewFromConfig(cfg)
	fakeInstanceID := "i-04b465ab13755d687"
	if err := startInstance(context.TODO(), fakeInstanceID, client, true); err != nil {
		if !strings.Contains(err.Error(), "api error DryRunOperation: Request would have succeeded, but DryRun flag is set") {
			t.Errorf("fail to stop the instance: %v", err)
		}
		t.Skip()
	}

}

func TestCheckExpectedTime(t *testing.T) {
	t.Logf("test trigger time against working time")

	expectedOK := true
	actualTime := time.Now()
	zone, err := time.LoadLocation("")
	if err != nil {
		t.Errorf("configure timezone failed %v", err)
	}
	// t := time.Now()

	actualTime = time.Date(actualTime.Year(), actualTime.Month(), actualTime.Day(), 0, 0, 0, 0, zone)
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
	actualTime := time.Now()
	zone, err := time.LoadLocation("")
	if err != nil {
		t.Errorf("configure timezone failed %v", err)
	}
	// t := time.Now()
	actualTime = time.Date(actualTime.Year(), actualTime.Month(), actualTime.Day(), 0, 0, 0, 0, zone).Add(12 * time.Hour)

	working, err := checkExpectedTime(context.TODO(), actualTime)
	if err != nil {
		t.Errorf("format time failed: %v", err)
	}

	if working != expectedOK {
		t.Errorf("wrong time check current time is %s, while shows in work", actualTime.String())
	}
}
