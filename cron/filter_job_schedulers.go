package cron

// FilterJobSchedulers filters job schedulers.
func FilterJobSchedulers(schedulers []*JobScheduler, predicate func(*JobScheduler) bool) []*JobScheduler {
	var output []*JobScheduler
	for _, js := range schedulers {
		if predicate(js) {
			output = append(output, js)
		}
	}
	return output
}
