package main

import (
	"context"
	"log"
	"net/http"

	"google.golang.org/grpc/health"
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
