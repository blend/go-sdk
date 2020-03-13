package grpcutil_test

import (
	"context"
	"testing"

	"github.com/blend/go-sdk/assert"
	v1 "github.com/blend/go-sdk/grpcutil/grpcutil_test/v1"
)

func numbers(values ...float64) []*v1.Number {
	output := make([]*v1.Value, len(values))
	for index := range values {
		output[index] = &v1.Number{Value: values[index]}
	}
	return output
}

func Test_CalculatorServer_Add(t *testing.T) {
	assert := assert.New(t)

	server := new(CalculatorServer)
	res, err := server.Add(context.TODO(), &v1.Numbers{
		Values: numbers(1, 2, 3, 4),
	})
}
