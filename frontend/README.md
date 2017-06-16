# Frontend

Frontend implements the ping server and depends on two backend servers.

## Usage

```
frontend -h
```

```
Usage of frontend:
  -bar string
    	The bar service address
  -foo string
    	The foo service address
  -grpc string
    	The gRPC listen address (default "127.0.0.1:8080")
  -health string
    	The health listen address (default "127.0.0.1:8008")
  -http string
    	The HTTP listen address (default "127.0.0.1:80")
  -region string
    	The compute region
```
