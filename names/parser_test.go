/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package names

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestNames(t *testing.T) {
	assert := assert.New(t)

	names := map[string]Name{}
	names["John Doe"] = Name{"", "John", "", "Doe", ""}
	names["Mr Anthony R Von Fange III"] = Name{"Mr.", "Anthony", "R", "Von Fange", "III"}
	names["Sara Ann Fraser"] = Name{"", "Sara Ann", "", "Fraser", ""}
	names["Adam"] = Name{"", "Adam", "", "", ""}
	names["Jonathan Smith"] = Name{"", "Jonathan", "", "Smith", ""}
	names["Anthony R Von Fange III"] = Name{"", "Anthony", "R", "Von Fange", "III"}
	names["Anthony Von Fange III"] = Name{"", "Anthony", "", "Von Fange", "III"}
	names["Mr John Doe"] = Name{"Mr.", "John", "", "Doe", ""}
	names["Justin White Phd"] = Name{"", "Justin", "", "White", "PhD"}
	names["Mark P Williams"] = Name{"", "Mark", "P", "Williams", ""}
	// Preserves the case of compound last name words
	names["Aaron bin Omar"] = Name{"", "Aaron", "", "bin Omar", ""}
	names["AARON BIN OMAR"] = Name{"", "Aaron", "", "Bin Omar", ""}
	names["Aaron ibn Omar"] = Name{"", "Aaron", "", "ibn Omar", ""}
	// Title case all caps
	names["MICHAEL J SPOLARICH JR"] = Name{"", "Michael", "J", "Spolarich", "Jr"}
	// Title case lowers
	names["michael j spolarich jr"] = Name{"", "Michael", "J", "Spolarich", "Jr"}
	// Preserve camel if provided
	names["LaTosha McMichaels"] = Name{"", "LaTosha", "", "McMichaels", ""}
	// Can't camel if we get consistent case
	names["LATOSHA MCMICHAELS"] = Name{"", "Latosha", "", "Mcmichaels", ""}
	names[""] = Name{"", "", "", "", ""}
	names["Dr"] = Name{"Dr.", "", "", "", ""}

	for rawName, expectedResult := range names {
		result := Parse(rawName)
		assert.Equal(expectedResult, result)
	}
}
