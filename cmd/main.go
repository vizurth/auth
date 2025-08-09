package main

import (
	"github.com/vizurth/auth/internal/server"
	desc "github.com/vizurth/auth/pkg/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

const grpcPort = ":50051"

func main() {
	userServer := server.NewServer()

	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)

	desc.RegisterUserServer(s, userServer)

	log.Printf("Starting gRPC server on port %s", grpcPort)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
