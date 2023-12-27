/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package pagerduty

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/r2"
)

// GetService gets a service
func (hc HTTPClient) GetService(ctx context.Context, id string) (service Service, err error) {
	var res *http.Response
	res, err = hc.Request(ctx,
		r2.OptGet(),
		r2.OptPath(fmt.Sprintf("/services/%s", id)),
	).Do()
	if err != nil {
		return
	}
	if statusCode := res.StatusCode; statusCode < 200 || statusCode > 299 {
		statusErr := ErrNon200Status
		if statusCode == 404 {
			statusErr = Err404Status
		}
		err = ex.New(statusErr, ex.OptMessagef("method: post, path: /services/%s, status: %d", id, statusCode))
		return
	}
	defer res.Body.Close()
	var output map[string]Service
	if err = json.NewDecoder(res.Body).Decode(&output); err != nil {
		err = ex.New(err)
		return
	}

	service, ok := output["service"]
	if !ok {
		err = ex.New("JSON response did not include the service field")
	}

	return
}
