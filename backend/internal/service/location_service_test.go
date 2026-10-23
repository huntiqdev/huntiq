package service

import (
	"context"
	"testing"

	"github.com/your-org/used-tools-marketplace/internal/models"
)

func TestLocationServiceCoordinateRangeValidation(t *testing.T) {
	svc := &LocationService{}

	testCases := []struct {
		name    string
		lat     float64
		lng     float64
		wantErr bool
	}{
		{"Valid Addis Ababa", 9.03, 38.74, false},
		{"Valid Hawassa", 7.05, 38.47, false},
		{"Invalid Lat High", 95.0, 38.74, true},
		{"Invalid Lat Low", -91.0, 38.74, true},
		{"Invalid Lng High", 9.03, 185.0, true},
		{"Invalid Lng Low", 9.03, -185.0, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := svc.ReverseGeocode(context.Background(), tc.lat, tc.lng)
			if tc.wantErr {
				if err == nil {
					t.Errorf("Expected error for coordinates (%f, %f), but got nil", tc.lat, tc.lng)
				}
			} else {
				// We expect "geocoder service not initialized" since geocoder is nil,
				// NOT "coordinates out of valid range"
				if err != nil && err.Error() == "coordinates out of valid range (lat: [-90, 90], lng: [-180, 180])" {
					t.Errorf("Unexpected coordinate range error for valid coordinates (%f, %f)", tc.lat, tc.lng)
				}
			}
		})
	}
}

func TestLocationServiceValidation(t *testing.T) {
	svc := &LocationService{}

	// Test missing name_en
	loc := &models.LocationHierarchy{
		NameEn:       "",
		LocationType: "woreda",
	}
	_, err := svc.CreateLocation(context.Background(), loc)
	if err == nil {
		t.Error("Expected error when creating location without name_en")
	}

	// Test invalid ID on update
	updateLoc := &models.LocationHierarchy{
		ID:     0,
		NameEn: "Test",
	}
	err = svc.UpdateLocation(context.Background(), updateLoc)
	if err == nil {
		t.Error("Expected error when updating location with ID <= 0")
	}
}
