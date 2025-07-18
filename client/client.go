package main

import (
	"context"
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gocarina/gocsv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"strconv"
	"sync"
	dp "telemetry/datapoint"
	pb "telemetry/proto"
	"time"
)

/*type DataPoint struct {
	BikeID       string  `csv:"Bike_id" json:"Bike_id"`
	Timestamp    int64   `csv:"Timestamp" json:"Timestamp"`
	Latitude     float64 `csv:"Latitude" json:"Latitude"`
	Longitude    float64 `csv:"Longitude" json:"Longitude"`
	BatteryLevel float64 `csv:"Battery_level" json:"Battery_level"`
	Charging     bool    `csv:"Charging" json:"Charging"`
}*/

func Test() {
	file, err := os.OpenFile("fleet_sim/bike_data0b8dc573.csv", os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var dataPoints []dp.DataPoint

	if err := gocsv.UnmarshalFile(file, &dataPoints); err != nil {
		panic(err)
	}

	for _, dp := range dataPoints {
		fmt.Printf("Parsed: %+v\n", dp)
	}
}

// Implement readDataPointsFromCSV for DataPoint
func readDataPointsFromCSV(bikeID string) []dp.DataPoint {
	//print current directory
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting working directory:", err)
	}
	fmt.Println("Current working directory:", dir)

	filename := "fleet_sim/bike_data" + bikeID + ".csv"
	file, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var dataPoints []dp.DataPoint

	if err := gocsv.UnmarshalFile(file, &dataPoints); err != nil {
		panic(err)
	}
	return dataPoints
}

/*type CurrentDataPoint struct {
	dp DataPoint
	mu sync.Mutex // sync access...
}

func (data_p *CurrentDataPoint) GetCurrent() DataPoint {
	data_p.mu.Lock()
	defer data_p.mu.Unlock()
	return data_p.dp
}

func (data_p *CurrentDataPoint) SetCurrentDataPoint(dp DataPoint) {
	data_p.mu.Lock()
	data_p.dp = dp
	data_p.mu.Unlock()
}*/

// returns time difference...
func getTimeDiff(t1Str, t2Str string) time.Duration {
	t1, err := time.Parse(time.RFC3339, t1Str)
	if err != nil {
		panic(err)
	}

	t2, err := time.Parse(time.RFC3339, t2Str)
	if err != nil {
		panic(err)
	}

	return t1.Sub(t2) // returns time.Duration

}

func startBikeSimulation(dataPoints []dp.DataPoint, c_data_point *dp.CurrentDataPoint) {

	// init with first
	c_data_point.SetCurrentDataPoint(dataPoints[0])
	// run forever
	counter := 1
	waittime := dataPoints[1].Timestamp - dataPoints[0].Timestamp // just sets a time in between new readings from bike
	for {
		// pause for waittime_in_milliseconds
		time.Sleep(time.Duration(waittime) * time.Millisecond)
		index := counter % len(dataPoints) // just loop over the datapoints forever
		c_data_point.SetCurrentDataPoint(dataPoints[index])
		counter++
	}
}

func runGRPCClient(cp *dp.CurrentDataPoint, ctx context.Context) {
	cc, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer cc.Close()

	if err != nil {
		log.Fatalf("Client could not connect: %v", err)
	}

	client := pb.NewTelemetryServiceClient(cc)

	//ctx := context.Background()

	// run simulation in thread...

	stream, err := client.ChargingStatus(ctx)
	if err != nil {
		log.Fatalf("could not start streaming: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			datapoint := cp.GetCurrentDataPoint()
			req := &pb.ChargingStatusRequest{
				BikeId:       datapoint.BikeID,
				BatteryLevel: float32(datapoint.BatteryLevel),
				IsCharging:   datapoint.Charging,
				Timestamp:    strconv.FormatInt(datapoint.Timestamp, 10),
				Lat:          datapoint.Latitude,
				Long:         datapoint.Longitude,
			}
			if err := stream.Send(req); err != nil {
				log.Fatalf("failed to send message: %v", err)
			}
			log.Printf("Sent charging status update")
			time.Sleep(time.Millisecond * 1000) // send data every second
		}
	}

}

func runMQTTPublisher(cp *dp.CurrentDataPoint, ctx context.Context) {
	broker := "tcp://localhost:1883" // or your broker URL
	opts := mqtt.NewClientOptions().AddBroker(broker)
	opts.SetClientID("bike1234")

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal("Connection error:", token.Error())
	}
	log.Println("Connected to MQTT broker")

	for {
		select {
		case <-ctx.Done():
			log.Println("[mqtt] Shutting down...")
			return
		default:
			datap := cp.GetCurrentDataPoint()
			payload, err := json.Marshal(datap)
			if err != nil {
				log.Println("Failed to marshal:", err)
				break
			}
			//topic := dp.BikeID // this is what we should listen for...
			token := client.Publish("bike-pw/"+datap.BikeID, 0, true, payload)
			token.Wait()
			if err := token.Error(); err != nil {
				log.Println("Publish error:", err)
			} else {
				log.Println("Published to topic: bike-pw-1/" + datap.BikeID)
			}
			time.Sleep(time.Second * 1)
		}
	}

}

func main() {
	//Test()
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run client.go bike device id")
		return
	}

	bike_id := os.Args[1]
	datapoints := readDataPointsFromCSV(bike_id)

	current_data_point := dp.NewCurrentDataPoint()
	go startBikeSimulation(datapoints, current_data_point)
	ctx, cancel := context.WithCancel(context.Background())

	wg := sync.WaitGroup{}
	wg.Add(3)
	/*go func() {
		defer wg.Done()
		runGRPCClient(current_data_point, ctx)
	}()

	go func() {
		defer wg.Done()
		runMQTTPublisher(current_data_point, ctx)
	}()*/

	go func() {
		defer wg.Done()
		runCoapServer1(current_data_point, ctx)
	}()

	go func() {
		time.Sleep(60 * time.Second) // run for half a minute then shut down...
		cancel()                     // signal all goroutines to stop
	}()

	wg.Wait()
	log.Println("client is done...both grpc and mqtt has stopped")

}
