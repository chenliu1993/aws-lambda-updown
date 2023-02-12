package pkg

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/cenkalti/backoff/v4"
)

func HandlerReq(ctx context.Context, req Request) error {
	instanceID := &req.InstanceID
	client, err := New(instanceID, req.ApiKey, req.ApiSecret)
	if err != nil {
		return err
	}

	state, err := checkInstanceStatus(ctx, *instanceID, client)
	if err != nil {
		return err
	}
	log.Println("the current state is ", state.Name)

	switch state.Name {
	case "running":
		ok, err := checkExpectedTime(ctx)
		if err != nil {
			return err
		}
		if !ok {
			//start
			stopInstance(ctx, *instanceID, client)
		}
	case "stopped":
		ok, err := checkExpectedTime(ctx)
		if err != nil {
			return err
		}
		if ok {
			//start
			startInstance(ctx, *instanceID, client)
			// startInstance(ctx, *instanceID, client)
		}
	default:
		log.Printf("instance is under wrong state: %s", state.Name)
	}
	return nil
}

func checkInstanceStatus(ctx context.Context, instanceID string, client *ec2.Client) (*ec2types.InstanceState, error) {
	log.Println("check the instance status ", instanceID)
	// Describe the status
	output, err := client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		InstanceIds: []string{
			instanceID,
		},
	})

	if err != nil {
		return nil, err
	}

	return output.Reservations[0].Instances[0].State, nil
}

func startInstance(ctx context.Context, instanceID string, client *ec2.Client) error {
	log.Println("start the instance ", instanceID)
	input := &ec2.StartInstancesInput{
		InstanceIds: []string{
			instanceID,
		},
		// DryRun: aws.Bool(true),
	}

	log.Printf("begin start instance")
	_, err := client.StartInstances(ctx, input)
	if err != nil {
		log.Printf("start met error %v", err.Error())
		return err
	}
	log.Printf("start instance op done")
	return backoff.Retry(func() error {
		state, err := checkInstanceStatus(ctx, instanceID, client)
		if err != nil {
			return err
		}
		if state.Name != "running" {
			return fmt.Errorf("the instance %s is still in %s state: ", instanceID, state.Name)
		}
		log.Println("started")
		return nil
	}, backoff.WithMaxRetries(backoff.NewConstantBackOff(1*time.Second), 3))
}

func stopInstance(ctx context.Context, instanceID string, client *ec2.Client) error {
	log.Println("stop the instance ", instanceID)
	input := &ec2.StopInstancesInput{
		InstanceIds: []string{
			instanceID,
		},
		// DryRun: aws.Bool(true),
	}
	_, err := client.StopInstances(ctx, input)
	if err != nil {
		log.Printf("stop met error %v", err.Error())
		return err
	}
	return backoff.Retry(func() error {
		state, err := checkInstanceStatus(ctx, instanceID, client)
		if err != nil {
			return err
		}
		if state.Name != "stopped" {
			return fmt.Errorf("the instance %s is still in %s state: ", instanceID, state.Name)
		}
		log.Println("stopped")
		return nil
	}, backoff.WithMaxRetries(backoff.NewConstantBackOff(1*time.Second), 3))
}

func checkExpectedTime(ctx context.Context) (bool, error) {
	// zone := time.FixedZone("Asia/Tokyo", 9*3600)
	zone, err := time.LoadLocation("")
	if err != nil {
		log.Printf("configure timezone failed %v", err)
		return false, err
	}
	t := time.Now()
	currentTime := t.Add(9 * time.Hour)
	startTime := time.Date(t.Year(), t.Month(), t.Day(), 9, 0, 0, 0, zone)
	endTime := time.Date(t.Year(), t.Month(), t.Day(), 18, 0, 0, 0, zone)
	log.Printf("current time is %s, start time is %s, end time is %s", currentTime.String(), startTime.String(), endTime.String())
	return currentTime.Sub(startTime) >= 0 && endTime.Sub(currentTime) > 0, nil
}
