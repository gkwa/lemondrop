package lemondrop

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/patrickmn/go-cache"
)

type RegionComponents struct {
	Region         string
	RegionDesc     string
	RegionCode     string
	RegionFriendly string
	City           string
}

type RegionDetails map[string]RegionComponents

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

func GetAllAwsRegions() (RegionDetails, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))
	if err != nil {
		return nil, err
	}

	svc := ssm.NewFromConfig(cfg)

	// copied from https://stackoverflow.com/a/72357524/1495086

	regionDetails := make(RegionDetails)
	var nextToken *string
	for {
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
			regionDetails[region] = RegionComponents{
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

func GetRegionDetails() (RegionDetails, error) {
	cachePath, err := getCachePath()
	if err != nil {
		return RegionDetails{}, err
	}

	regions, err := fetchFromCache()
	if err != nil {
		return RegionDetails{}, err
	}

	cacheHitDebugMsg := "regions in cache"

	if len(regions) != 0 {
		slog.Debug(cacheHitDebugMsg, "hit", true)
		return regions, nil
	}

	slog.Debug(cacheHitDebugMsg, "hit", false)

	regions, err = GetAllAwsRegions()
	if err != nil {
		return RegionDetails{}, err
	}

	jsonBytes, err := json.MarshalIndent(regions, "", "  ")
	if err != nil {
		return RegionDetails{}, err
	}
	regionsCache.Set(cacheKey, string(jsonBytes), cache.DefaultExpiration)
	regionsCache.SaveFile(cachePath)

	return regions, nil
}

func GetRegions(writer io.Writer) {
	regions, err := GetRegionDetails()
	if err != nil {
		panic(err)
	}

	for _, rDetail := range regions {
		fmt.Fprintf(writer, "%s\n", rDetail.Region)
	}
}
