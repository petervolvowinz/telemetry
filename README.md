# Basic telemetry service tutorial gRPC,MQTT, CoAp


3 different IOT protocols gRPC, MQTT and CoAp, using Golang as language. We currently have a multi-threaded
example application and a very basic telemetry example:

```
...
message ChargingStatusRequest{
  string bike_id      = 1;
  float battery_level = 2;
  bool is_charging    = 3;
  string timestamp    = 4;
  double lat = 5;
  double long = 6;
}
... 
```

Here expressed in protobuf syntax. 

**gRPC**

gRPC is a rpc protocol that lets you stream data bi-directional. It has great language support and is known for
its performance. Major drawback is that it is resource demanding and that it depends on http2 which is sometimes blocked
out from corporate walls - basically because don't support HTTP2.

**MQTT**

MQTT runs a broker service where clients can publish or subscribe to messages.
To install a local running MQTT broker to run local development tests against we can use
For Mac:
```
$ brew install mosquitto
$ brew services start mosquitto

```

This lets us run local testing and the default port for the local broker is:


```
    broker := "tcp://localhost:1883" // public test broker
	clientID := "bike123"
```

In MQTT we publish data through a concept of topics: 
Topics are defined as a tree structure with slash (/) as delimiters. Strings are UTF-8 encoded, they act as message 
channels which are then available for client subscription. Topics are also case-sensitive, below a MQTT topic example:

```
  topic : "bike-telemetry/123456"
```

A client can choose to listen to a specific message or it can use wildcards:
```
    "bike-telemetry/123456"   subscribe to messages matching this string 
    "bike-telemetry/#"        subscribe to any message published under bike-telemetry
```



```Go
broker := "tcp://localhost:1883" // public test broker
	clientID := "bike123"
	// Define topic to listen to (e.g., a specific bike ID or wildcard)
	topic := "bike-pw/#" // any bike id
	if len(os.Args) > 1 {
		topic = os.Args[1] // topic as program arg
	}

	// Create MQTT client options
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID) // sets the id of this client.

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


```



**CoAp**

Another lightweight RESTful IOT protocol that is mostly used in resource constrained microcontrollers...<TODO>

**gRPC implementation**

gRPC uses a interface defintion language approach where you from a protobuf file generates code for the specific language.

We define our telemetry service::
```protobuf

service TelemetryService {
  rpc ChargingStatus (stream ChargingStatusRequest) returns (ChargingStatusResponse);
}

message ChargingStatusRequest{
  string bike_id      = 1;
  float battery_level = 2;
  bool is_charging    = 3;
  string timestamp    = 4;
  double lat = 5;
  double long = 6;
}

message ChargingStatusResponse{
  string status = 1;
}
```


**MQTT implementation**
**CoAp implementation**
**Open street map web sever**


