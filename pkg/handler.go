package pkg

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/cenkalti/backoff/v4"
)

func HandlerReq(ctx context.Context, req Request) error {
	instanceID := &req.InstanceID
	client, err := New(instanceID)
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
		// stop
		stopInstance(ctx, *instanceID, client)
		return nil
	case "stopped":
		//start
		startInstance(ctx, *instanceID, client)
	default:
		return fmt.Errorf("instance is under wrong state: %s", state.Name)
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
		DryRun: aws.Bool(true),
	}

	output, err := client.StartInstances(ctx, input)
	if err != nil {
		return err
	}
	return backoff.Retry(func() error {
		if output.StartingInstances[0].CurrentState.Name != "running" {
			return fmt.Errorf("the instance %s is still in %s state: ", *output.StartingInstances[0].InstanceId, output.StartingInstances[0].CurrentState.Name)
		}
		return nil
	}, backoff.WithMaxRetries(backoff.NewConstantBackOff(200*time.Millisecond), 3))
}

func stopInstance(ctx context.Context, instanceID string, client *ec2.Client) error {
	log.Println("stop the instance ", instanceID)
	input := &ec2.StopInstancesInput{
		InstanceIds: []string{
			instanceID,
		},
		DryRun: aws.Bool(true),
	}
	output, err := client.StopInstances(ctx, input)
	if err != nil {
		return err
	}
	return backoff.Retry(func() error {
		if output.StoppingInstances[0].CurrentState.Name != "stopped" {
			return fmt.Errorf("the instance %s is still in %s state: ", *output.StoppingInstances[0].InstanceId, output.StoppingInstances[0].CurrentState.Name)
		}
		return nil
	}, backoff.WithMaxRetries(backoff.NewConstantBackOff(200*time.Millisecond), 3))
}
