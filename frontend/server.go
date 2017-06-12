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
	"google.golang.org/grpc/peer"
)

type server struct {
	bar ping.PingClient
	foo ping.PingClient
}

func (s *server) Ping(ctx context.Context, in *ping.Request) (*ping.Response, error) {
	p, ok := peer.FromContext(ctx)
	if ok {
		log.Printf("Ping request from %s", p.Addr)
	}

	barResponse, err := s.bar.Ping(context.Background(), &ping.Request{})
	if err != nil {
		log.Printf("Error calling service bar: %v", err)
		return nil, err
	}

	fooResponse, err := s.foo.Ping(context.Background(), &ping.Request{})
	if err != nil {
		log.Printf("Error calling service foo: %v", err)
		return nil, err
	}

	log.Printf("bar version: %s", barResponse.Version)
	log.Printf("foo version: %s", fooResponse.Version)

	return &ping.Response{Message: "Pong", Version: version}, nil
}
