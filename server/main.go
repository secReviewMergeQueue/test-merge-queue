package main

import (
	"log"
	"net"

	pb "github.com/saumeya/train-ticketing/api/proto"
	bookingservice "github.com/saumeya/train-ticketing/server/bookservice"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen on port 50051: %v", err)
	}

	s := grpc.NewServer()

	// Register the booking service
	bookingService := bookingservice.NewBookingService()
	pb.RegisterBookingServiceServer(s, bookingService)

	log.Printf("Welcome to Train Ticketing Service %v \n", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
