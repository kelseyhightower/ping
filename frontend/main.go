// Copyright 2017 Google Inc. All Rights Reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
//
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/kelseyhightower/ping"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/reflection"
)

const (
	version = "v1.0.0"
)

var (
	listenAddr   string
	serviceBAddr string
	serviceCAddr string
)

type server struct {
	bc ping.PingClient
	cc ping.PingClient
}

func main() {
	flag.StringVar(&listenAddr, "listen-addr", "127.0.0.1:50051", "The gRPC listen address")
	flag.StringVar(&serviceBAddr, "service-b-addr", "127.0.0.1:50052", "The address for service B")
	flag.StringVar(&serviceCAddr, "service-c-addr", "127.0.0.1:50053", "The address for service C")
	flag.Parse()

	log.Println("Starting frontend service ...")
	log.Println("Listening on", listenAddr)

	// Create a gRPC client for service B.
	bconn, err := grpc.Dial(serviceBAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer bconn.Close()
	bc := ping.NewPingClient(bconn)

	// Create a gRPC client for service C.
	cconn, err := grpc.Dial(serviceCAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer cconn.Close()
	cc := ping.NewPingClient(cconn)

	s := grpc.NewServer()
	ping.RegisterPingServer(s, &server{bc, cc})
	reflection.Register(s)

	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		log.Fatal(s.Serve(ln))
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	log.Printf("Shutdown signal received shutting down gracefully...")

	s.GracefulStop()
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
