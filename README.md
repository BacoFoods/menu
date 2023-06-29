
# Setup environment

```bash
# MacOS
brew install protobuf

go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0

protoc --go_out=. --go-grpc_out=. proto/menu/*.proto
```