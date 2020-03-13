//go:generate protoc -I ./v1 --go_out=plugins=grpc:./v1 ./v1/calculator.proto

package grpcutil_test
