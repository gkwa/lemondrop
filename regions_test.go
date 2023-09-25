package lemondrop

import (
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

func TestGetRegionDetails(t *testing.T) {
	regions, err := GetRegionDetails()
	if err != nil {
		panic(err)
	}
	regionDetails := regions["us-west-2"]
	wantCity := "Oregon"
	gotCity := regionDetails.City

	if gotCity != wantCity {
		t.Errorf("got %q want %q", gotCity, wantCity)
	}
}
