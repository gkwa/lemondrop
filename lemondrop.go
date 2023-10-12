package lemondrop

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/patrickmn/go-cache"
	mymazda "github.com/taylormonacelli/forestfish/mymazda"
	"github.com/taylormonacelli/somespider"
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
	regions, err := fetchCachedRegions()
	if err != nil {
		return RegionDetails{}, err
	}

	if len(regions) != 0 {
		slog.Debug("regions in cache", "hit", true)
		return regions, nil
	}

	// cache miss
	slog.Debug("regions in cache", "hit", false)

	regions, err = GetAllAwsRegions()
	if err != nil {
		return RegionDetails{}, err
	}

	jsonBytes, err := json.MarshalIndent(regions, "", "  ")
	if err != nil {
		return RegionDetails{}, err
	}
	regionsCache.Set(cacheKey, string(jsonBytes), cache.DefaultExpiration)
	defer peristCacheToDisk()

	return regions, nil
}

func peristCacheToDisk() error {
	cachePath, err := somespider.GenPath(relCachPath)
	if err != nil {
		return err
	}

	// prepare to persist cache to disk:
	cacheSnapshot := regionsCache.Items()

	gob.Register(RegionDetails{})

	// serialize using gob:
	file, _ := os.Create(cachePath)
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(cacheSnapshot)
	if err != nil {
		slog.Error("encode", "error", err.Error())
		return err
	}
	defer file.Close()
	slog.Debug("checking existance of file cache", "exists", mymazda.FileExists(cachePath))
	return nil
}

func WriteRegions(writer io.Writer, showDesc bool) {
	regions, err := GetRegionDetails()
	if err != nil {
		panic(err)
	}

	for _, rDetail := range regions {
		if showDesc {
			fmt.Fprintf(writer, "%s [%s]\n", rDetail.RegionCode, rDetail.RegionDesc)
		} else {
			fmt.Fprintf(writer, "%s\n", rDetail.RegionCode)
		}
	}
}
