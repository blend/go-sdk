/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

/*
Package status provides helpers for building standardized sla status endpoints in web services.

It provides types to wrap existing actions to track the success and error history of those actions
and then report that history through status endpoints.

The package provides a `status.Controller` type that is created with the `status.NewController(...)` method:

	    statusController := status.New(
			status.OptCheck(
				"redis",
				redisConnection, // implements `status.Checker` for you already
			),
			status.OptCheck(
				"postgres",
				dbConnection, // implements `status.Checker` for you already
			),
		)
		...
		app.Register(statusController)

The app will now have `/status/sla` and `/status/details` endpoints registered.
*/
package status	// import "github.com/blend/go-sdk/status"
