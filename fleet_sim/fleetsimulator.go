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

func simulateBikeData(bikeID string, startLat, startLng float64, startTime time.Time, numPoints int) []DataPoint {
	data := make([]DataPoint, 0, numPoints)
	batteryLevel := 100.0

	for i := 0; i < numPoints; i++ {
		lat := startLat + rand.Float64()*0.002 - 0.001
		lng := startLng + rand.Float64()*0.002 - 0.001

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
