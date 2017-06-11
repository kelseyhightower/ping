package main

import (
	"log"

	"github.com/kelseyhightower/ping"
	"golang.org/x/net/context"
	"google.golang.org/grpc/peer"
)

type server struct{}

func (s *server) Ping(ctx context.Context, in *ping.Request) (*ping.Response, error) {
	p, ok := peer.FromContext(ctx)
	if ok {
		log.Printf("Ping request from %s", p.Addr)
	}
	return &ping.Response{Message: "Pong", Version: version}, nil
}
