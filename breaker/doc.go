/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

/*
Package breaker provides a circuitbraker mechanism for dealing with flaky or unreliable counterparties.

The algorithm used for the state machine is described by Microsoft https://docs.microsoft.com/en-us/previous-versions/msp-n-p/dn589784(v=pandp.10)

An example of using a circuit breaker for an http call might be:

	    b := breaker.New()
		phoneHome := b.Intercept(async.ActionerFunc(func(ctx context.Context, _ interface{}) (interface{}, error) {
			return http.DefaultClient.Do(http.NewRequestWithContext(ctx, http.VerbGet "https://google.com/robots.txt", nil))
	    })

In the above, `phoneHome` now will be wrapped with circuit breaker mechanics. You would call it with `phoneHome.Action(ctx, nil)` etc.
*/
package breaker	// import "github.com/blend/go-sdk/breaker"
