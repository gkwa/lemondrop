package lemondrop

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

const (
	cacheFilePath   = "regions_cache.json"
	cacheExpiration = 1 * 24 * time.Hour
)

type CachedRegions struct {
	Timestamp string         `json:"timestamp"`
	Regions   []types.Region `json:"regions"`
}

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

func readRegionsFromCache() ([]types.Region, error) {
	cacheFile, err := os.Open(cacheFilePath)
	if err != nil {
		return nil, err
	}
	defer cacheFile.Close()

	var cachedRegions CachedRegions
	err = json.NewDecoder(cacheFile).Decode(&cachedRegions)
	if err != nil {
		return nil, err
	}

	cachedTime, err := time.Parse(time.RFC3339, cachedRegions.Timestamp)
	if err != nil {
		return nil, err
	}

	// Check if the cached regions have expired
	if time.Since(cachedTime) > cacheExpiration {
		return nil, fmt.Errorf("cache has expired")
	}

	return cachedRegions.Regions, nil
}

func writeRegionsToCache(regions []types.Region) error {
	cacheFileDir := filepath.Dir(cacheFilePath)
	if err := os.MkdirAll(cacheFileDir, 0o755); err != nil {
		return err
	}

	cacheFile, err := os.Create(cacheFilePath)
	if err != nil {
		return err
	}
	defer cacheFile.Close()

	cachedRegions := CachedRegions{
		Timestamp: time.Now().Format(time.RFC3339),
		Regions:   regions,
	}

	encodedRegions, err := json.MarshalIndent(cachedRegions, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(cacheFilePath, encodedRegions, 0o644)
	if err != nil {
		return err
	}

	return nil
}
