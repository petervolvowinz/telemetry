package main

import (
	"context"
	"encoding/json"
	"github.com/plgd-dev/go-coap/v3/udp"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	dp "telemetry/datapoint"
	"time"
)

func main() {
	coapconnection, err := udp.Dial("localhost:5688")
	if err != nil {
		log.Fatalf("Error dialing: %v", err)
	}

	path := "/bike-pw/telemetry"
	if len(os.Args) > 1 { // you can add customized path as command line parameter
		path = os.Args[1]
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60) // run for a minute
	defer cancel()

	go func() {
		for {
			resp, err := coapconnection.Get(ctx, path)
			if err != nil {
				log.Fatalf("Error sending request: %v", err)
			}
			bytesread, err := io.ReadAll(resp.Body())
			if err != nil {
				log.Fatalf("failed to read response body: %v", err)
			}

			data := &dp.DataPoint{}
			err = json.Unmarshal(bytesread, data)
			if err != nil {
				log.Fatalf("Error unmarshalling response: %v", err)
			} else {
				log.Println("Bike id: %v", data.BikeID)
				log.Println("Battery level: %v", data.BatteryLevel)
				log.Println("Charging status: %v", data.Charging)
			}
			time.Sleep(time.Second)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Println("coap observer shutting down...")
}
