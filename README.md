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

Lightweight protocol that ... <TODO>

**CoAp**

Another lightweight RESTful IOT protocol that is mostly used in resource constrained microcontrollers...<TODO>

**gRPC implementation**
**MQTT implementation**
**CoAp implementation**

