package main

import (
	"fmt"
	"log"
	"net"
	"os"

	restaurantsSvc "github.com/cagodoy/tenpo-restaurants-api/rpc/restaurants"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/cagodoy/tenpo-challenge/lib/proto"
	_ "github.com/lib/pq"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5030"
		log.Println("missing env variable PORT, using default value...")
	}

	srv := grpc.NewServer()
	svc := restaurantsSvc.New()

	pb.RegisterRestaurantServiceServer(srv, svc)
	reflection.Register(srv)

	log.Println("Starting Restaurants service...")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		log.Fatalf("Failed to list: %v", err)
	}

	log.Println(fmt.Sprintf("Restaurants service, Listening on: %v", port))

	if err := srv.Serve(lis); err != nil {
		log.Fatalf("Fatal to serve: %v", err)
	}
}
