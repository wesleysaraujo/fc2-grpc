## go-grpc

## gRPC

### Gerar Protocol Buffer

`protoc --proto_path=proto proto/*.proto --go_out=pb`

### Gerar as Stubs
`protoc --proto_path=proto proto/*.proto --go_out=pb --go-grpc_out=pb`

### Rodar o server gRPC
`go run cmd/server/server.go`

### Rodar o Cliente
`go run cmd/client/client.go`