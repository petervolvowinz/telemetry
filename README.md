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

This lets us run local testing . There are also url based test brokers to test against.




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


