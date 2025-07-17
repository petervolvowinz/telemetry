package main

import (
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	pb "telemetry/proto"
)

type server struct {
	pb.UnimplementedTelemetryServiceServer
}

func (s *server) ChargingStatus(stream pb.TelemetryService_ChargingStatusServer) error {
	log.Println("ChargingStatus stream started")
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// Client finished sending
			return stream.SendAndClose(&pb.ChargingStatusResponse{
				Status: "All charging data received",
			})
		}
		if err != nil {
			log.Printf("Error receiving from stream: %v", err)
			return err
		}
		log.Printf("Received charging status: bike_id=%s, battery_level=%f, is_charging=%v, timestamp=%s",
			req.BikeId, req.BatteryLevel, req.IsCharging, req.Timestamp)
		// Here you could process/store the telemetry data
	}
}

func (s *server) mustEmbedUnimplementedTelemetryServiceServer() {
	log.Printf("must embedd telemetry service server")
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterTelemetryServiceServer(s, &server{})
	log.Println("telemetry server listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Server failure: %v", err)
	}
}
