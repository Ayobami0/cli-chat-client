# CLI Chat
A command line chat app build with golang and gRPC. Connects to [server](https://github.com/Ayobami0/cli-chat-server)

## Install
### Using the executable (linux only)
```
cd cmd/
./cli-chat
# optionally add cmd to you path to run anywhere
```
### Build from source
__Requirements__
- [golang](https://golang.org/doc/install)
- make
- [protoc](https://grpc.io/docs/protoc-installation/)

1. Install golang plugin for the [protoc compiler](https://grpc.io/docs/languages/go/quickstart/#prerequisites)
```
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
```
2. Export the SERVER_ADDR envrionment variable
```
export SERVER_ADDR=<server_address>
```
> [!IMPORTANT]
> If the variable is not set the binary won't work
3. Run make
```
# You may optionally supply a output name by setting the OUT environment variable
# export OUT=<custom_name>
make build
```
4. Run generated binary
```
./cli-chat
```
