all: generate-protos

init:
	@go get -u github.com/golang/protobuf/protoc-gen-go

generate:
	@go generate github.com/blend/go-sdk/examples/grpcutil/full/...
