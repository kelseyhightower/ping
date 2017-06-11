package main

import (
	"context"
	"log"

	"github.com/kelseyhightower/ping"
	"google.golang.org/grpc/peer"
)

type server struct {
	bc ping.PingClient
	cc ping.PingClient
}

func (s *server) Ping(ctx context.Context, in *ping.Request) (*ping.Response, error) {
	p, ok := peer.FromContext(ctx)
	if ok {
		log.Printf("Ping request from %s", p.Addr)
	}

	rb, err := s.bc.Ping(context.Background(), &ping.Request{})
	if err != nil {
		log.Printf("Error calling service B: %v", err)
		return nil, err
	}

	rc, err := s.cc.Ping(context.Background(), &ping.Request{})
	if err != nil {
		log.Printf("Error calling service C: %v", err)
		return nil, err
	}

	log.Printf("Service B version: %s", rb.Version)
	log.Printf("Service C version: %s", rc.Version)

	return &ping.Response{Message: "Pong", Version: version}, nil
}
