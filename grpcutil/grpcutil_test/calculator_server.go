package grpcutil_test

import (
	"context"
	"io"

	v1 "github.com/blend/go-sdk/grpcutil/grpcutil_test/v1"
)

var (
	_ v1.CalculatorServer = (*CalculatorServer)(nil)
)

// CalculatorServer is the server for the calculator.
type CalculatorServer struct{}

// Add adds a fixed set of numbers.
func (CalculatorServer) Add(_ context.Context, values *v1.Numbers) (*v1.Number, error) {
	var output float64
	for _, value := range values.Values {
		output += value
	}
	return &v1.Number{
		Value: output,
	}, nil
}

// AddStream adds a stream of numbers.
func (CalculatorServer) AddStream(stream v1.Calculator_AddStreamServer) error {
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
func (CalculatorServer) Subtract(_ context.Context, values *v1.Numbers) (*v1.Number, error) {
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
func (CalculatorServer) SubtractStream(stream v1.Calculator_SubtractStreamServer) error {
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
func (CalculatorServer) Multiply(_ context.Context, values *v1.Numbers) (*v1.Number, error) {
	var output float64
	for _, value := range values.Values {
		output *= value
	}
	return &v1.Number{
		Value: output,
	}, nil
}

// MultiplyStream multiplies a stream of numbers.
func (CalculatorServer) MultiplyStream(stream v1.Calculator_MultiplyStreamServer) error {
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

		output *= number.Value
	}
}

// Divide divides a fixed set of numbers.
func (CalculatorServer) Divide(_ context.Context, values *v1.Numbers) (*v1.Number, error) {
	if len(values.Values) == 0 {
		return nil, nil
	}
	output := values.Values[0]
	for _, value := range values.Values[1:] {
		output = output / value
	}
	return &v1.Number{
		Value: output,
	}, nil
}

// DivideStream divides a stream of numbers.
func (CalculatorServer) DivideStream(stream v1.Calculator_DivideStreamServer) error {
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
			output = output / number.Value
		}
	}
}
