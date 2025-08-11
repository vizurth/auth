package main

import (
	"context"
	"fmt"
	"github.com/vizurth/auth/internal/config"
	"github.com/vizurth/auth/internal/postgres"
	"github.com/vizurth/auth/internal/server"
	desc "github.com/vizurth/auth/pkg/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	ctx := context.Background()

	cfg, _ := config.NewConfig()

	db, err := postgres.New(ctx, cfg.Postgres)

	if err != nil {
		log.Fatal(err)
	}

	userServer := server.NewServer(db)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)

	desc.RegisterUserServer(s, userServer)

	log.Printf("Starting gRPC server on port %s", fmt.Sprintf(":%s", cfg.Port))
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
