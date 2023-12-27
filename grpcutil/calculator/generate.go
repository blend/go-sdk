/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

//go:generate protoc -I ./v1 --go_out=./v1 --go-grpc_out=./v1 ./v1/calculator.proto

package calculator
