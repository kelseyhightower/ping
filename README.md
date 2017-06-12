# Ping

This repository holds example microservice composed of the following components:

* frontend - depends on the bar and foo services
* bar - a microservice that implements the ping server
* foo - a microservice that implements the ping server
* client - gRPC client that talks to the frontend

## Generate gRPC code

```
protoc -I ./ ./ping.proto --go_out=plugins=grpc:.
```
