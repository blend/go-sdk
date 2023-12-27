/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package calculator

import (
	"context"
	"io"

	v1 "github.com/blend/go-sdk/grpcutil/calculator/v1"
)

// Server is the server for the calculator.
type Server struct {
	v1.CalculatorServer
}

// Add adds a fixed set of numbers.
func (Server) Add(_ context.Context, values *v1.Numbers) (*v1.Number, error) {
	var output float64
	for _, value := range values.Values {
		output += value
	}
	return &v1.Number{
		Value: output,
	}, nil
}

// AddStream adds a stream of numbers.
func (Server) AddStream(stream v1.Calculator_AddStreamServer) error {
	var output float64
	var number *v1.Number
	var err error
	for {
		select {
		case <-stream.Context().Done():
			return stream.SendAndClose(&v1.Number{
				Value: output,
			})
		default:
		}

		number, err = stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&v1.Number{
				Value: output,
			})
		}

		output += number.Value
	}
}

// Subtract subtracts a fixed set of numbers.
func (Server) Subtract(_ context.Context, values *v1.Numbers) (*v1.Number, error) {
	if len(values.Values) == 0 {
		return nil, nil
	}
	output := values.Values[0]
	for _, value := range values.Values[1:] {
		output -= value
	}
	return &v1.Number{
		Value: output,
	}, nil
}

// SubtractStream subtracts a stream of numbers.
func (Server) SubtractStream(stream v1.Calculator_SubtractStreamServer) error {
	var outputSet bool
	var output float64
	var number *v1.Number
	var err error
	for {
		select {
		case <-stream.Context().Done():
			return stream.SendAndClose(&v1.Number{
				Value: output,
			})
		default:
		}

		number, err = stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&v1.Number{
				Value: output,
			})
		}
		if !outputSet {
			output = number.Value
			outputSet = true
		} else {
			output -= number.Value
		}
	}
}

// Multiply multiplies a fixed set of numbers.
func (Server) Multiply(_ context.Context, values *v1.Numbers) (*v1.Number, error) {
	if len(values.Values) == 0 {
		return nil, nil
	}
	output := values.Values[0]
	for _, value := range values.Values[1:] {
		output *= value
	}
	return &v1.Number{
		Value: output,
	}, nil
}

// MultiplyStream multiplies a stream of numbers.
func (Server) MultiplyStream(stream v1.Calculator_MultiplyStreamServer) error {
	var output float64
	var outputSet bool
	var number *v1.Number
	var err error
	for {
		select {
		case <-stream.Context().Done():
			return stream.SendAndClose(&v1.Number{
				Value: output,
			})
		default:
		}

		number, err = stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&v1.Number{
				Value: output,
			})
		}

		if !outputSet {
			output = number.Value
			outputSet = true
		} else {
			output *= number.Value
		}
	}
}

// Divide divides a fixed set of numbers.
func (Server) Divide(_ context.Context, values *v1.Numbers) (*v1.Number, error) {
	if len(values.Values) == 0 {
		return nil, nil
	}
	output := values.Values[0]
	for _, value := range values.Values[1:] {
		output /= value
	}
	return &v1.Number{
		Value: output,
	}, nil
}

// DivideStream divides a stream of numbers.
func (Server) DivideStream(stream v1.Calculator_DivideStreamServer) error {
	var outputSet bool
	var output float64
	var number *v1.Number
	var err error
	for {
		select {
		case <-stream.Context().Done():
			return stream.SendAndClose(&v1.Number{
				Value: output,
			})
		default:
		}

		number, err = stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&v1.Number{
				Value: output,
			})
		}
		if !outputSet {
			output = number.Value
			outputSet = true
		} else {
			output /= number.Value
		}
	}
}
