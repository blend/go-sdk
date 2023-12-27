/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package consistenthash

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/blend/go-sdk/assert"
)

func spew(v interface{}) string {
	output, _ := json.Marshal(v)
	return string(output)
}

func assignmentCounts(assignments map[string][]string) map[string]int {
	output := make(map[string]int)
	for bucket, items := range assignments {
		output[bucket] = len(items)
	}
	return output
}

// setsAreDisjoint asserts that a given list of sets has no overlapping items with any of the other sets.
//
// this is helpful to prove that there are no duplicated items within sets.
func setsAreDisjoint(its *assert.Assertions, itemSets ...[]string) {
	for thisSetIndex, thisSet := range itemSets {
		for _, thisItem := range thisSet {
			for otherSetIndex, otherSet := range itemSets {
				if thisSetIndex != otherSetIndex {
					for _, otherItem := range otherSet {
						its.NotEqual(thisItem, otherItem)
					}
				}
			}
		}
	}
}

func setExistsInOtherSets(its *assert.Assertions, items []string, otherSets ...[]string) {
	for _, item := range items {
		its.True(findItem(item, otherSets...), fmt.Sprintf("Item was not found in any supplied assignments: %s", item))
	}
}

func setsAreEvenlySized(its *assert.Assertions, itemCount int, sets ...[]string) {
	// assert that the item count is represented in the sets somewhat evenly
	setCount := len(sets)
	itemsPerSet := itemCount / setCount
	minSize := 0
	maxSize := 2 * itemsPerSet

	for _, set := range sets {
		its.True(len(set) > minSize && len(set) < maxSize, fmt.Sprintf("sets should be between %d and %d in size, actual: %d", minSize, maxSize, len(set)))
	}
}

// setIsSorted asserts the set, and the set modified by `sort.Strings` are the same.
func setIsSorted(its *assert.Assertions, set []string) {
	its.NotEmpty(set)
	trialSorted := make([]string, len(set))
	copy(trialSorted, set)
	sort.Strings(trialSorted)
	its.Equal(trialSorted, set)
}

// itsConsistent asserts that a given before and after set differ by a given maximum delta.
func itsConsistent(its *assert.Assertions, old, new []string, maxDelta int) {
	var delta int
	for _, oldItem := range old {
		if !findItem(oldItem, new) {
			delta++
		}
	}
	its.True(delta < maxDelta, fmt.Sprintf("The item arrays should differ only by %d items (differed by %d items)\n\told: %s\n\tnew: %s\n", maxDelta, delta, strings.Join(old, ", "), strings.Join(new, ", ")))
}

func findItem(item string, assignments ...[]string) bool {
	for _, assignment := range assignments {
		for _, assignmentItem := range assignment {
			if item == assignmentItem {
				return true
			}
		}
	}
	return false
}

func assignmentsAreEqual(its *assert.Assertions, oldAssignments, newAssignments map[string][]string) {
	if len(oldAssignments) != len(newAssignments) {
		its.Fail(fmt.Sprintf("assignments are not equal; length of maps differ: %d vs. %d", len(oldAssignments), len(newAssignments)))
	}

	for oldBucketName, oldBucket := range oldAssignments {
		newBucket, ok := newAssignments[oldBucketName]
		if !ok {
			its.Fail(fmt.Sprintf("assignments are not equal; old bucket %s not found in new assignments", oldBucketName))
		}
		if len(oldBucket) != len(newBucket) {
			its.Fail(fmt.Sprintf("assignments are not equal; length of assignment lists for %s differ: %d vs. %d", oldBucketName, len(oldBucket), len(newBucket)))
		}
		for _, oldBucketItem := range oldBucket {
			if !findItem(oldBucketItem, newBucket) {
				its.Fail(fmt.Sprintf("assignments are not equal; old bucket %s item %s not found in new bucket", oldBucketName, oldBucketItem))
			}
		}
	}
}
