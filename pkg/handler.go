package pkg

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/cenkalti/backoff/v4"
)

var layout = "2006-01-02 15:04:05"

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
		ok, err := checkExpectedTime(ctx, req.StartHour, req.StopHour, req.Timezone)
		if err != nil {
			return err
		}
		if !ok {
			//start
			stopInstance(ctx, *instanceID, client)
		}
	case "stopped":
		ok, err := checkExpectedTime(ctx, req.StartHour, req.StopHour, req.Timezone)
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

func checkExpectedTime(ctx context.Context, start_time, stop_time, timezone string) (bool, error) {
	// default to UTC
	zone, err := time.LoadLocation("")
	if timezone != "" {
		zone, err = time.LoadLocation(timezone)
	}
	if err != nil {
		log.Printf("configure timezone failed %v", err)
		return false, err
	}
	t := time.Now().In(zone)

	// currentTime := t.Add(9 * time.Hour)
	var (
		startTime  = time.Date(t.Year(), t.Month(), t.Day(), 7, 0, 0, 0, zone)
		endTime    = time.Date(t.Year(), t.Month(), t.Day(), 22, 0, 0, 0, zone)
		start, end int
	)
	if start_time != "" {
		start, err = strconv.Atoi(start_time)
		if err != nil {
			log.Printf("parse start time failed %v/%s", err, start_time)
			return false, err
		}
		startTime = time.Date(t.Year(), t.Month(), t.Day(), start, 0, 0, 0, zone)
	}

	if stop_time != "" {
		end, err = strconv.Atoi(stop_time)
		if err != nil {
			log.Printf("parse start time failed %v/%s", err, stop_time)
			return false, err
		}
		endTime = time.Date(t.Year(), t.Month(), t.Day(), end, 0, 0, 0, zone)
	}

	log.Printf("current time is %s, start time is %s, end time is %s", t.String(), startTime.String(), endTime.String())
	return t.Sub(startTime) >= 0 && endTime.Sub(t) > 0, nil
}
