# gRPC

## tools
```sh

brew install protoc

brew install protoc-gen-go
brew install protoc-gen-go-grpc

brew install buf
brew install grpcurl

go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2

```
- ENSURE ALL the above tools are in your $PATH, esp those installed via "go install"
- important folders are internal/api and internal/api_proto
  - internal/api_proto -> contains protobuf contract and its autogen interface class files
  - internal/api/grpc_server.go -> contains the RPC server implementation(gRPC API)
  
## basics
- create .proto file and place it in the respective version folder, eg. /internal/api_proto/v1
- ensure the buf tool is installed on the dev machine

```sh
cd internal/api_proto/v1 #ensure you are in this as the current folder
touch main.proto
touch buf.yaml
touch buf.gen.yaml
```

- Add relevant contents into the buf.yaml file, which declares the major dps and core buf config
- add relevant contents into the buf.gen.yaml file, which controls destination of protoc-generated outputs
- run the buf update for to download any modules that are declared as dependencies in the buf.yaml file, which produces a buf.lock file

```sh
buf dep update
```

## working
- generate the interface classes for the RPC server, which will be generated according to the contents declared in the proto file.
- this command reads the *.proto files in the current directory
- for each x.proto, it outputs x.pb.go and x_grpc.pb.go files
```sh
buf generate
```

## implementation
- implement the interface files in the internal/api/grpc_server.go file by importing the v1 interface files
- register the implememntation with the grpc mux @ /internal/api/servers.go+myGrpcServer:Run
- register any extra dependencies with the grpc mux if necessary(e.g. if you are using the grpc-gateway plugin etc)
- run the service and test the endpoint logic

## gRPC gateway
- Ensure grpc gateway implementation is registered in the grpc mux instantiation
- Ensure HTTP endpoints to serve the gRPC-redirected requests is registered with http mux

## testing
```sh

grpcurl --plaintext 127.0.0.1:8082 describe

grpcurl -d '{"userName":"JohnCena"}' --plaintext 127.0.0.1:8082  mybackend.v1.MybackendGrpcSvc/GetUserInfoByUsername

grpcurl -d '{"userId":"SOME-UUID"}' --plaintext 127.0.0.1:8082  mybackend.v1.MybackendGrpcSvc/GetUserInfoById

curl -X 'GET' \
  'http://localhost:8081/v1/user/byusername?userName=zack' \
  -H 'accept: application/json' \
  -H 'Authorization: test'

curl -X 'GET' \
  'http://localhost:8081/v1/user/byid' \
  -H 'accept: application/json' \
  -H 'Authorization: test'


```

- Hooray!!

