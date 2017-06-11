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
	"log"
	"net/http"

	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type healthHandler struct {
	healthServer *health.Server
}

func httpHealthServer(server *health.Server) http.Handler {
	return &healthHandler{server}
}

func (h *healthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hcr, err := h.healthServer.Check(context.Background(), &healthpb.HealthCheckRequest{""})
	if err != nil {
		log.Println("Error checking gRPC server health", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	switch hcr.Status.String() {
	case "UNKNOWN":
		w.WriteHeader(http.StatusServiceUnavailable)
	case "SERVING":
		w.WriteHeader(http.StatusOK)
	case "NOT_SERVING":
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}
