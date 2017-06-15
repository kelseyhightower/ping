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
	"encoding/json"
	"log"
	"net/http"

	"github.com/kelseyhightower/ping"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type pingHandler struct {
	localAddr string
}

func httpPingServer(addr string) http.Handler {
	return &pingHandler{addr}
}

type httpResponse struct {
	BarVersion string `json:"bar_version"`
	FooVersion string `json:"foo_version"`
	Hostname   string `json:"hostname"`
	Message    string `json:"message"`
	Region     string `json:"region"`
	Version    string `json:"version"`
}

func (p *pingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := grpc.Dial(p.localAddr, grpc.WithInsecure())
	if err != nil {
		log.Println("Error calling the local ping server", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	defer conn.Close()
	client := ping.NewPingClient(conn)

	md := metadata.New(map[string]string{})
	ctx := context.Background()
	grpcResponse, err := client.Ping(ctx, &ping.Request{}, grpc.Trailer(&md))
	if err != nil {
		log.Println("Error calling the local ping server", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	response := httpResponse{
		BarVersion: md["barversion"][0],
		FooVersion: md["fooversion"][0],
		Hostname:   md["hostname"][0],
		Message:    grpcResponse.Message,
		Region:     md["region"][0],
		Version:    md["version"][0],
	}

	data, err := json.MarshalIndent(&response, "", "  ")
	if err != nil {
		log.Println("Error marshalling HTTP response:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
	return
}
