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
	"log"

	"github.com/kelseyhightower/ping"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type server struct {
	bar      ping.PingClient
	foo      ping.PingClient
	hostname string
	region   string
	version  string
}

func (s *server) Ping(ctx context.Context, in *ping.Request) (*ping.Response, error) {
	// Call the bar service and extract the version from the response metadata.
	barMetadata := metadata.New(map[string]string{})
	barCtx := context.Background()
	_, err := s.bar.Ping(barCtx, &ping.Request{}, grpc.Trailer(&barMetadata))
	if err != nil {
		log.Printf("Error calling bar service: %v", err)
		return nil, err
	}

	barVersion := barMetadata["version"][0]

	// Call the foo service and extract the version from the response metadata.
	fooMetadata := metadata.New(map[string]string{})
	fooCtx := context.Background()
	_, err = s.foo.Ping(fooCtx, &ping.Request{}, grpc.Trailer(&fooMetadata))
	if err != nil {
		log.Printf("Error calling foo service: %v", err)
		return nil, err
	}

	fooVersion := fooMetadata["version"][0]

	// Set the reponse metadata that will be send back to the client.
	md := metadata.New(map[string]string{
		"barVersion": barVersion,
		"fooVersion": fooVersion,
		"hostname":   s.hostname,
		"region":     s.region,
		"version":    s.version,
	})

	if err := grpc.SetTrailer(ctx, md); err != nil {
		log.Printf("Error setting the response metadata: %v", err)
	}

	return &ping.Response{Message: "pong"}, nil
}
