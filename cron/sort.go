/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package cron

// JobSchedulersByJobNameAsc is a wrapper that sorts job schedulers
// by the job name ascending.
type JobSchedulersByJobNameAsc []*JobScheduler

// Len implements sorter.
func (s JobSchedulersByJobNameAsc) Len() int {
	return len(s)
}

// Swap implements sorter.
func (s JobSchedulersByJobNameAsc) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less implements sorter.
func (s JobSchedulersByJobNameAsc) Less(i, j int) bool {
	return s[i].Name() < s[j].Name()
}
