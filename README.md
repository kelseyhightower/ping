# Ping

## Generate gRPC code

```
protoc -I ./ ./ping.proto --go_out=plugins=grpc:.
```
