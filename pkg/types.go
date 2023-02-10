package pkg

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

// only for initialize some of the functions of aws client
type Caller interface {
	DescribeInstances(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)
	StartInstances(ctx context.Context, params *ec2.StartInstancesInput, optFns ...func(*ec2.Options)) (*ec2.StartInstancesOutput, error)
	StopInstances(ctx context.Context, params *ec2.StopInstancesInput, optFns ...func(*ec2.Options)) (*ec2.StopInstancesOutput, error)
}

type Request struct {
	InstanceID string `json:"instance_id"`
}
