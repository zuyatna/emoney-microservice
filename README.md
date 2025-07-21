# emoney-microservice
A mini microservices-based digital wallet and payment service powered by gRPC, Postgres, Redis, and Elasticsearch

## Preparation Tools
### Protobuf
Install protobuf https://protobuf.dev/installation/

### Google api annotations
Navigate to where you store your projects
`cd /path/to/your/projects/`

Clone the repository
git clone https://github.com/googleapis/googleapis

### Go plugin gRPC
- `go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28`
- `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2`
