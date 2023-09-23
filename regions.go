package lemondrop

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type DetailedRegion struct {
	Region         string
	RegionDesc     string
	RegionCode     string
	RegionFriendly string
	City           string
}

func getCity(str string) (string, string, error) {
	pattern := `([^(]+) (\(([^)]+)\))?`

	re := regexp.MustCompile(pattern)

	submatches := re.FindStringSubmatch(str)

	if len(submatches) >= 2 {
		regionFriendly := strings.Trim(submatches[1], " ")
		city := strings.Trim(submatches[3], " ")
		return regionFriendly, city, nil
	}
	return "", "", nil
}

func GetAllAwsRegions() (map[string]DetailedRegion, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))
	if err != nil {
		return nil, err
	}

	// Create an SSM client
	svc := ssm.NewFromConfig(cfg)

	regionDetails := make(map[string]DetailedRegion)
	var nextToken *string
	for {
		// Request all regions, paginating the results if needed
		input := &ssm.GetParametersByPathInput{
			Path:      aws.String("/aws/service/global-infrastructure/regions"),
			NextToken: nextToken,
		}

		resp, err := svc.GetParametersByPath(context.TODO(), input)
		if err != nil {
			return nil, err
		}

		// For each region, get the "longName" for the region
		for _, parameter := range resp.Parameters {
			region := (*parameter.Name)[strings.LastIndex(*parameter.Name, "/")+1:]

			regionInfo, err := svc.GetParameter(context.TODO(), &ssm.GetParameterInput{
				Name: aws.String("/aws/service/global-infrastructure/regions/" + region + "/longName"),
			})
			if err != nil {
				return nil, err
			}

			regionDesc := *regionInfo.Parameter.Value
			rf, city, err := getCity(regionDesc)
			if err != nil {
				panic(err)
			}
			regionDetails[region] = DetailedRegion{
				City:           city,
				Region:         region,
				RegionCode:     region,
				RegionDesc:     regionDesc,
				RegionFriendly: rf,
			}
		}

		// Pull in the next page of regions if needed
		nextToken = resp.NextToken
		if nextToken == nil {
			break
		}
	}
	return regionDetails, nil
}

func GetAllAwsRegions1() ([]types.Region, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
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
