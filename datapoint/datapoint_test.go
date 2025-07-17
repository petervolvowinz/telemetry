package datapoint

import (
	"testing"
	"time"
)

func TestCurrentDataPoint_SetAndGet(t *testing.T) {
	cp := NewCurrentDataPoint()

	expected := DataPoint{
		BikeID:       "bike-test-01",
		Timestamp:    time.Now().Unix(),
		Latitude:     51.509865,
		Longitude:    -0.118092,
		BatteryLevel: 78.5,
		Charging:     true,
	}

	// Set the value
	cp.SetCurrentDataPoint(expected)

	// Small sleep to allow goroutine to process
	time.Sleep(10 * time.Millisecond)

	// Get the value
	actual := cp.GetCurrentDataPoint()

	// Validate fields
	if actual.BikeID != expected.BikeID {
		t.Errorf("BikeID mismatch: got %s, expected %s", actual.BikeID, expected.BikeID)
	}
	if actual.Timestamp != expected.Timestamp {
		t.Errorf("Timestamp mismatch: got %d, expected %d", actual.Timestamp, expected.Timestamp)
	}
	if actual.Latitude != expected.Latitude {
		t.Errorf("Latitude mismatch: got %f, expected %f", actual.Latitude, expected.Latitude)
	}
	if actual.Longitude != expected.Longitude {
		t.Errorf("Longitude mismatch: got %f, expected %f", actual.Longitude, expected.Longitude)
	}
	if actual.BatteryLevel != expected.BatteryLevel {
		t.Errorf("BatteryLevel mismatch: got %f, expected %f", actual.BatteryLevel, expected.BatteryLevel)
	}
	if actual.Charging != expected.Charging {
		t.Errorf("Charging mismatch: got %v, expected %v", actual.Charging, expected.Charging)
	}
}
