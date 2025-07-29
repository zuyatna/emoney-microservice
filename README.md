# emoney-microservice
A mini microservices-based digital wallet and payment service powered by gRPC, Postgres, Redis, and Elasticsearch

## Preparation Tools
### Protobuf
Install protobuf https://protobuf.dev/installation/

### Go plugin gRPC
- `go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28`
- `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2`

### Google api annotations
Navigate to where you store your projects
`cd /path/to/your/projects/`

Clone the repository
git clone https://github.com/googleapis/googleapis

### Adding google annotation
- `go get google.golang.org/genproto/googleapis/api/annotations`

### Generate proto to protobuffer
```
protoc --proto_path=proto \
       --proto_path=../../googleapis \
       --go_out=server/pb --go_opt=paths=source_relative \
       --go-grpc_out=server/pb --go-grpc_opt=paths=source_relative \
       --grpc-gateway_out=server/pb --grpc-gateway_opt=paths=source_relative \
       proto/account.proto
```
> Assuming that you clone google api in outer root project.
>
> Note: Change `--proto-path=../../googleapis`, based on your project structures

### Setup Postgres, Redis, RabbitMQ, and Elasticsearch in Docker
- PostgreSQL
```
docker run --name my-postgres \
       -e POSTGRES_PASSWORD=passwordrahasia \
       -p 5432:5432 \
       -v postgres-data:/var/lib/postgresql/data \
       -d postgres
```

Access postgres via docker:
`docker exec -it my-postgres psql -U postgres`

Get all databases:
`\l`

Creating new database:
`CREATE DATABASE emoney;`

Access database:
`\c {database_name}`

- Redis
```
docker run --name my-redis \
       -p 6379:6379 \
       -v redis-data:/data \
       -d redis
```

- RabbitMQ
```
docker run -d --hostname my-rabbit --name some-rabbit \
    -e RABBITMQ_DEFAULT_USER=user \
    -e RABBITMQ_DEFAULT_PASS=password \
    -p 5672:5672 \
    -p 15672:15672 \
    rabbitmq:4-management
```
- Elasticsearch
```
docker network create elastic

docker pull docker.elastic.co/elasticsearch/elasticsearch:9.0.4

docker run --name es01 --net elastic -p 9200:9200 -it -m 1GB docker.elastic.co/elasticsearch/elasticsearch:9.0.4
```
