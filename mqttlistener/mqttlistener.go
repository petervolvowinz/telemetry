package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	broker := "tcp://localhost:1883" // public test broker
	clientID := "bike123"

	// Define topic to listen to (e.g., a specific bike ID or wildcard)
	topic := "bike-pw/#" // any bike id

	// Create MQTT client options
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)

	// Create and connect the client
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal("Failed to connect:", token.Error())
	}
	log.Println("Connected to MQTT broker")

	if token := client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("Received message on [%s]: %s\n", msg.Topic(), msg.Payload())
	}); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	log.Printf("Subscribed to topic: %s\n", topic)

	// Wait for Ctrl+C to exit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Println("Disconnecting...")
	client.Disconnect(250)
	log.Println("Listener shut down")
}
