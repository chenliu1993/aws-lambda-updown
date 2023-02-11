package pkg

import (
	"context"
	"errors"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

func New(instanceID *string, api_key, api_secret string) (*ec2.Client, error) {
	log.Println("initialize client for ", *instanceID)
	if *instanceID == "" {
		log.Println("You must supply an instance ID (-i INSTANCE-ID")
		return nil, errors.New("please input a instance ID")
	}

	// Default profile is OK
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	provider := credentials.NewStaticCredentialsProvider(api_key, api_secret, "")
	cfg.Credentials = aws.NewCredentialsCache(provider)

	client := ec2.NewFromConfig(cfg)
	return client, nil
}
