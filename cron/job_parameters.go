package cron

import (
	"context"
	"encoding/json"
	"net/url"
)

// JobParametersFromForm creates a parameter values set from url values.
func JobParametersFromForm(formValues url.Values) JobParameters {
	output := make(JobParameters)
	for key, values := range formValues {
		if len(values) == 0 {
			output[key] = ""
			continue
		}
		output[key] = values[0]
	}
	return output
}

// JobParameterValuesFromJSON creates a parameter values set from json data.
func JobParameterValuesFromJSON(data []byte) (JobParameters, error) {
	output := make(JobParameters)
	if err := json.Unmarshal(data, &output); err != nil {
		return nil, err
	}
	return output, nil
}

// JobParameters is a loose association to map[string]string.
type JobParameters = map[string]string

type contextKeyJobParameters struct{}

// WithJobParameters adds job invocation parameter values to a context.
func WithJobParameters(ctx context.Context, values JobParameters) context.Context {
	return context.WithValue(ctx, contextKeyJobParameters{}, values)
}

// GetJobParameters gets parameter values from a given context.
func GetJobParameters(ctx context.Context) JobParameters {
	if value := ctx.Value(contextKeyJobParameters{}); value != nil {
		if typed, ok := value.(JobParameters); ok {
			return typed
		}
	}
	return nil
}
