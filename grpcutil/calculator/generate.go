/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

//go:generate protoc -I ./v1 --go_out=./v1 --go-grpc_out=./v1 ./v1/calculator.proto

package calculator
