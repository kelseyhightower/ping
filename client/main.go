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
	"fmt"
	"log"

	"github.com/kelseyhightower/ping"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	version = "v1.0.0"
)

var (
	serverAddr string
)

func main() {
	flag.StringVar(&serverAddr, "server", "127.0.0.1:50051", "The ping server address")
	flag.Parse()

	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := ping.NewPingClient(conn)

	response, err := c.Ping(context.Background(), &ping.Request{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(response.Message)
}
