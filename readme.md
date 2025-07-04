# Demo
This repo holds a simple proto handshake between a go server and a go client

# Instructions

Generates message.pb.go
- `protoc --go_out=proto proto/message.proto`
Build
- `go build -o bin/server.exe ./cmd/server`
- `go build -o bin/client.exe ./cmd/client`
Run
- `./bin/server.exe`
- `./bin/client.exe`