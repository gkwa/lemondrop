package lemondrop

import (
	"fmt"
	"testing"
)

func TestGetCity(t *testing.T) {
	str := "Europe (Spain)"
	wantCity := "Spain"
	wantFriendlyRegion := "Europe"
	gotFriendlyRegion, gotCity, _ := getCity(str)

	if gotFriendlyRegion != wantFriendlyRegion {
		t.Errorf("got %q want %q", gotFriendlyRegion, wantFriendlyRegion)
	}

	if gotCity != wantCity {
		t.Errorf("got %q want %q", gotCity, wantCity)
	}
}

func TestGetAllAwsRegions(t *testing.T) {
	regionsMap, err := GetAllAwsRegions()
	if err != nil {
		panic(err)
	}
	for region := range regionsMap {
		x := regionsMap[region]
		fmt.Printf("%s: region: %s, city: %s\n", x.RegionCode, x.RegionFriendly, x.City)
	}
}
