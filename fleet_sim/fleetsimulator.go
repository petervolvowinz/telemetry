package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"telemetry/genid"
	"time"
)

type DataPoint struct {
	BikeID       string
	Timestamp    int64
	Latitude     float64
	Longitude    float64
	BatteryLevel float64
	Charging     bool
}

var a_san_fransisco_route = []struct {
	Lat float64
	Lng float64
}{
	{37.7955, -122.3937}, // Ferry Building
	{37.8000, -122.3980}, // Embarcadero
	{37.8076, -122.4200}, // Fort Mason
	{37.8090, -122.4270}, // Marina Blvd
	{37.8105, -122.4360}, // Crissy Field
	{37.8077, -122.4477}, // Golden Gate Bridge View
	{37.7989, -122.4591}, // Park Presidio Blvd
	{37.7715, -122.4687}, // 14th Ave & Fulton
	{37.7694, -122.4862}, // Golden Gate Park (west side)
	{37.7691, -122.5107}, // Great Highway
	{37.7603, -122.5108}, // Ocean Beach
}

func simulateBikeData(bikeID string, startLat, startLng float64, startTime time.Time, numPoints int) []DataPoint {
	data := make([]DataPoint, 0, numPoints)
	batteryLevel := 100.0

	route_start := rand.Intn(len(a_san_fransisco_route))
	for i := 0; i < numPoints; i++ {

		lat := a_san_fransisco_route[route_start].Lat
		lng := a_san_fransisco_route[route_start].Lng

		route_start = (route_start + 1) % len(a_san_fransisco_route)

		var charging bool
		if i%15 == 0 && i != 0 {
			batteryLevel += 5.0 + rand.Float64()*5.0
			if batteryLevel > 100.0 {
				batteryLevel = 100.0
			}
			charging = true
		} else {
			batteryLevel -= 0.5 + rand.Float64()*1.5
			if batteryLevel < 0 {
				batteryLevel = 0
			}
			charging = false
		}

		t := startTime.Add(time.Duration(i) * 500 * time.Millisecond)
		data = append(data, DataPoint{
			BikeID:       bikeID,
			Timestamp:    t.UnixMilli(),
			Latitude:     lat,
			Longitude:    lng,
			BatteryLevel: batteryLevel,
			Charging:     charging,
		})
	}

	return data
}

func saveDataPerBike(data []DataPoint) {
	filename := "bike_data" + data[0].BikeID + ".csv"
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// define csv writer stuff
	writer := csv.NewWriter(file)
	defer writer.Flush()
	// write csv header
	writer.Write([]string{"Bike_id", "Timestamp", "Latitude", "Longitude", "Battery_level", "Charging"})
	for _, datapoint := range data {
		datarow := []string{
			datapoint.BikeID,
			strconv.FormatInt(int64(datapoint.Timestamp), 10),
			strconv.FormatFloat(datapoint.Latitude, 'f', -1, 64),
			strconv.FormatFloat(datapoint.Longitude, 'f', -1, 64),
			strconv.FormatFloat(datapoint.BatteryLevel, 'f', -1, 64),
			strconv.FormatBool(datapoint.Charging),
		}
		writer.Write(datarow)
	}
}

func main() {
	const numBikes = 50
	const pointsPerBike = 50
	// san fran
	startLat := 37.7749
	startLng := -122.4194
	startTime := time.Now()

	allData := make([]DataPoint, 0, numBikes*pointsPerBike)

	// Simulate data for each bike
	bid := &genid.Bikeid{}
	for i := 1; i <= numBikes; i++ {
		bikeID := bid.GenBikeId(genid.ShortIdGenerator)
		data := simulateBikeData(bikeID, startLat, startLng, startTime, pointsPerBike)
		saveDataPerBike(data)
		allData = append(allData, data...)
	}

	// Write all data to CSV
	file, err := os.Create("ebike_fleet_simulated.csv")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	writer.Write([]string{"bike_id", "timestamp", "latitude", "longitude", "battery_level", "charging"})

	// Write data
	for _, dp := range allData {
		row := []string{
			dp.BikeID,
			strconv.FormatInt(int64(dp.Timestamp), 10),
			fmt.Sprintf("%.6f", dp.Latitude),
			fmt.Sprintf("%.6f", dp.Longitude),
			fmt.Sprintf("%.2f", dp.BatteryLevel),
			strconv.FormatBool(dp.Charging),
		}
		writer.Write(row)
	}

	fmt.Printf("✅ CSV file 'ebike_fleet_simulated.csv' created with %d datapoints.\n", numBikes*pointsPerBike)
}
