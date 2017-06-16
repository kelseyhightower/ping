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

	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/reflection"
)

const (
	version = "v1"
)

var (
	barAddr    string
	fooAddr    string
	grpcAddr   string
	healthAddr string
	httpAddr   string
	region     string
)

func main() {
	flag.StringVar(&barAddr, "bar", "", "The bar service address")
	flag.StringVar(&fooAddr, "foo", "", "The foo service address")
	flag.StringVar(&grpcAddr, "grpc", "127.0.0.1:8080", "The gRPC listen address")
	flag.StringVar(&healthAddr, "health", "127.0.0.1:8008", "The health listen address")
	flag.StringVar(&httpAddr, "http", "127.0.0.1:80", "The HTTP listen address")
	flag.StringVar(&region, "region", "", "The compute region")
	flag.Parse()

	log.Println("Starting frontend service ...")
	log.Printf("gRPC server listening on: %s", grpcAddr)
	log.Printf("Health server listening on: %s", healthAddr)
	log.Printf("HTTP server listening on: %s", httpAddr)

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal("Error getting hostname:", err)
	}

	// Create a gRPC client for service bar.
	barConn, err := grpc.Dial(barAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer barConn.Close()
	barClient := ping.NewPingClient(barConn)

	// Create a gRPC client for service foo.
	fooConn, err := grpc.Dial(fooAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer fooConn.Close()
	fooClient := ping.NewPingClient(fooConn)

	// Setup the gRPC server.
	grpcServer := grpc.NewServer()
	s := &server{barClient, fooClient, hostname, region, version}
	ping.RegisterPingServer(grpcServer, s)
	reflection.Register(grpcServer)

	grpcHealthServer := health.NewServer()
	grpcHealthServer.SetServingStatus("ping.Ping", 0)
	healthpb.RegisterHealthServer(grpcServer, grpcHealthServer)

	ln, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatal(err)
	}

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

	// Setup a HTTP server to proxy the gRPC server.
	mux := http.NewServeMux()
	mux.Handle("/ping", httpPingServer(grpcAddr))
	httpServer := http.Server{Addr: httpAddr, Handler: mux}

	go func() {
		log.Fatal(httpServer.ListenAndServe())
	}()

	grpcHealthServer.SetServingStatus("ping.Ping", 1)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	log.Printf("Shutdown signal received shutting down gracefully...")

	grpcServer.GracefulStop()
	healthServer.Shutdown(context.Background())
	httpServer.Shutdown(context.Background())
}
