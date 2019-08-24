package dbmetrics

import "github.com/blend/go-sdk/db"

// Metric and tag names etc.
const (
	MetricNameDBQuery        string = string(db.QueryFlag)
	MetricNameDBQueryElapsed string = MetricNameDBQuery + ".elapsed"

	TagQuery    string = "query"
	TagEngine   string = "engine"
	TagDatabase string = "database"
)
