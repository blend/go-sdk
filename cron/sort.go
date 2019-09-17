package cron

// JobSchedulerStatusesByJobNameAsc is a wrapper that sorts job schedulers by the job name ascending.
type JobSchedulerStatusesByJobNameAsc []JobSchedulerStatus

// Len implements sorter.
func (s JobSchedulerStatusesByJobNameAsc) Len() int {
	return len(s)
}

// Swap implements sorter.
func (s JobSchedulerStatusesByJobNameAsc) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less implements sorter.
func (s JobSchedulerStatusesByJobNameAsc) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}
