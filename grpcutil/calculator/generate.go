//go:generate protoc -I ./v1 --go_out=./v1 --go-grpc_out=./v1 ./v1/calculator.proto

package calculator
