package main

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/kelseyhightower/ping"
	"golang.org/x/net/context"
)

type server struct {
	hostname string
	region   string
	version  string
}

func (s *server) Ping(ctx context.Context, in *ping.Request) (*ping.Response, error) {
	// Set the reponse metadata that will be send back to the client.
	md := metadata.New(map[string]string{
		"hostname": s.hostname,
		"region":   s.region,
		"version":  s.version,
	})

	if err := grpc.SetTrailer(ctx, md); err != nil {
		log.Printf("Error setting the response metadata: %v", err)
	}

	return &ping.Response{Message: "pong"}, nil
}
