package lemondrop

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type cacheEntry struct {
	Regions []types.Region
	Expiry  time.Time
}

var cache map[string]cacheEntry
var cacheMutex sync.Mutex

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

func init() {
	cache = make(map[string]cacheEntry)
}

func GetAllAwsRegions() ([]types.Region, error) {
	region := "us-west-2" // fixme: arbitrary and add more for failover

	// Check if the regions are already cached
	cacheMutex.Lock()
	entry, found := cache[region]
	cacheMutex.Unlock()

	if found && time.Now().Before(entry.Expiry) {
		// Return the cached regions
		return entry.Regions, nil
	}

	// Regions not found in cache, fetch them from the API
	client, err := GetEc2Client(region)
	if err != nil {
		return nil, err
	}

	// Get a list of all AWS regions
	resp, err := client.DescribeRegions(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to describe AWS regions: %v", err)
	}

	regions := resp.Regions

	// Cache the regions
	cacheMutex.Lock()
	cache[region] = cacheEntry{
		Regions: regions,
		Expiry:  time.Now().Add(time.Hour), // Cache expiry time
	}
	cacheMutex.Unlock()

	return regions, nil
}
