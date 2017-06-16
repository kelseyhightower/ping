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
	"context"
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
	version = "v2"
)

var (
	grpcAddr   string
	healthAddr string
	region     string
)

func main() {
	flag.StringVar(&grpcAddr, "grpc", "127.0.0.1:8080", "The gRPC listen address")
	flag.StringVar(&healthAddr, "health", "127.0.0.1:8008", "The health listen address")
	flag.StringVar(&region, "region", "", "The compute region")
	flag.Parse()

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal("Error getting hostname:", err)
	}

	log.Println("Starting backend service ...")
	log.Printf("gRPC server listening on: %s", grpcAddr)
	log.Printf("Health server listening on: %s", healthAddr)

	ln, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	ping.RegisterPingServer(grpcServer, &server{hostname, region, version})
	reflection.Register(grpcServer)

	grpcHealthServer := health.NewServer()
	grpcHealthServer.SetServingStatus("ping.Ping", 0)
	healthpb.RegisterHealthServer(grpcServer, grpcHealthServer)

	go func() {
		log.Fatal(grpcServer.Serve(ln))
	}()

	// Setup a HTTP server for health checks.
	healthMux := http.NewServeMux()
	healthMux.Handle("/health", httpHealthServer(grpcHealthServer))
	healthServer := http.Server{Addr: healthAddr, Handler: healthMux}

	go func() {
		log.Fatal(healthServer.ListenAndServe())
	}()

	grpcHealthServer.SetServingStatus("ping.Ping", 1)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	log.Printf("Shutdown signal received shutting down gracefully...")

	grpcServer.GracefulStop()
	healthServer.Shutdown(context.Background())
}
