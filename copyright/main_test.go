/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package copyright

import (
	"testing"

	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/testutil"
)

func TestMain(m *testing.M) {
	testutil.MarkUpdateGoldenFlag()

	testutil.New(
		m,
		testutil.OptLog(logger.All()),
	).Run()
}

const goBuildTags1 = `//go:build tag1
// +build tag1

package main

import (
	"fmt"
)

func main() {
	fmt.Println("foo")
}
`
const goBuildTags2 = `// +build tag5
//go:build tag1 && tag2 && tag3
// +build tag1,tag2,tag3
// +build tag6

package main

import (
	"fmt"
)

func main() {
	fmt.Println("foo")
}

// +bulid tag9000
`

const goBuildTags3 = `//go:build tag1 & tag2

package main

import (
	"fmt"
)

func main() {
	fmt.Println("foo")
}

//go:build tag3
`

const goldenGoBuildTags1 = `//go:build tag1
// +build tag1

/*

Copyright (c) 2001 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package main

import (
	"fmt"
)

func main() {
	fmt.Println("foo")
}
`

const goldenGoBuildTags2 = `// +build tag5
//go:build tag1 && tag2 && tag3
// +build tag1,tag2,tag3
// +build tag6

/*

Copyright (c) 2001 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package main

import (
	"fmt"
)

func main() {
	fmt.Println("foo")
}

// +bulid tag9000
`

const goldenGoBuildTags3 = `//go:build tag1 & tag2

/*

Copyright (c) 2001 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package main

import (
	"fmt"
)

func main() {
	fmt.Println("foo")
}

//go:build tag3
`

const goldenTsReferenceTags = `/// <reference path="../types/testing.d.ts" />
/// <reference path="../types/something.d.ts" />
/// <reference path="../types/somethingElse.d.ts" />
/// <reference path="../types/somethingMore.d.ts" />
/// <reference path="../types/somethingLess.d.ts" />
/**
 * Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
 * Blend Confidential - Restricted
 */
export * from '../types/goodOnes'
`

const tsReferenceTags = `/// <reference path="../types/testing.d.ts" />
/// <reference path="../types/something.d.ts" />
/// <reference path="../types/somethingElse.d.ts" />
/// <reference path="../types/somethingMore.d.ts" />
/// <reference path="../types/somethingLess.d.ts" />
export * from '../types/goodOnes'
`

const goldenTsReferenceTag = `/// <reference path="../types/testing.d.ts" />
/**
 * Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
 * Blend Confidential - Restricted
 */

export * from '../types/goodOnes'
`

const tsReferenceTag = `/// <reference path="../types/testing.d.ts" />

export * from '../types/goodOnes'
`

const tsTest = `export * from '../types/goodOnes'
`

const goldenTs = `/**
 * Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
 * Blend Confidential - Restricted
 */
export * from '../types/goodOnes'
`
