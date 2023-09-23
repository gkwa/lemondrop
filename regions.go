package lemondrop

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func GetAllAwsRegions() ([]types.Region, error) {
	region := "us-west-2" // fixme: arbitrary and add more for failover

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, err
	}
	client := ec2.NewFromConfig(cfg)
	if err != nil {
		return nil, err
	}

	resp, err := client.DescribeRegions(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to describe AWS regions: %v", err)
	}

	return resp.Regions, nil
}
