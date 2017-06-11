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
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/kelseyhightower/ping"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

const (
	version = "v2.0.0"
)

var (
	grpcAddr string
	httpAddr string
)

func main() {
	flag.StringVar(&grpcAddr, "grpc", "127.0.0.1:50051", "The gRPC listen address")
	flag.StringVar(&httpAddr, "http", "127.0.0.1:80", "The HTTP listen address")
	flag.Parse()

	log.Println("Starting backend service ...")
	log.Printf("gRPC server listening on: %s", grpcAddr)
	log.Printf("HTTP server listening on: %s", httpAddr)

	ln, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	ping.RegisterPingServer(s, &server{})
	reflection.Register(s)

	healthServer := health.NewServer()
	healthServer.SetServingStatus("grpc.health.v1.helloservice", 0)
	healthpb.RegisterHealthServer(s, healthServer)

	go func() {
		log.Fatal(s.Serve(ln))
	}()

	healthServer.SetServingStatus("grpc.health.v1.helloservice", 1)

	http.Handle("/healthz", httpHealthServer(healthServer))
	go func() {
		log.Fatal(http.ListenAndServe(httpAddr, nil))
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	log.Printf("Shutdown signal received shutting down gracefully...")

	s.GracefulStop()
}
