package lemondrop

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func CreateConfig(region string) (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return aws.Config{}, err
	}
	return cfg, nil
}

func GetEc2Client(region string) (*ec2.Client, error) {
	config, err := CreateConfig(region)
	if err != nil {
		return nil, err
	}
	// Create an EC2 client
	return ec2.NewFromConfig(config), nil
}

func GetAllAwsRegions() ([]types.Region, error) {
	// Return cached regions if available
	if cachedRegions, err := readRegionsFromCache(); err == nil {
		return cachedRegions, nil
	}

	region := "us-west-2" // fixme: arbitrary and add more for failover

	client, err := GetEc2Client(region)
	if err != nil {
		panic(err)
	}

	// Get a list of all AWS regions
	resp, err := client.DescribeRegions(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to describe AWS regions: %v", err)
	}

	regions := resp.Regions

	// Cache the regions to disk
	if err := writeRegionsToCache(regions); err != nil {
		fmt.Printf("Warning: Failed to write regions cache to disk: %v\n", err)
	}

	return regions, nil
}
