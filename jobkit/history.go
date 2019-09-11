package jobkit

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/blend/go-sdk/cron"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/stringutil"
)

// PersistHistoryJSON writes the history to disk as a json file
func PersistHistoryJSON(job cron.Job, config JobConfig) func(ctx context.Context, log []cron.JobInvocation) error {
	return func(ctx context.Context, log []cron.JobInvocation) error {
		if !config.HistoryPersistedOrDefault() {
			return nil
		}
		historyDirectory := config.HistoryPathOrDefault()
		if _, err := os.Stat(historyDirectory); err != nil {
			if err := os.MkdirAll(historyDirectory, 0755); err != nil {
				return ex.New(err)
			}
		}
		historyPath := filepath.Join(historyDirectory, stringutil.Slugify(job.Name())+".json")
		f, err := os.Create(historyPath)
		if err != nil {
			return err
		}
		defer f.Close()
		return json.NewEncoder(f).Encode(log)
	}
}

// RestoreHistory restores history from disc.
func (jh JSONHistory) RestoreHistory(ctx context.Context) (output []cron.JobInvocation, err error) {
	if !jh.HistoryEnabledOrDefault() {
		return nil, nil
	}
	if !jh.HistoryPersistedOrDefault() {
		return nil, nil
	}
	historyPath := filepath.Join(jh.HistoryPathOrDefault(), stringutil.Slugify(jh.Name)+".json")
	if _, statErr := os.Stat(historyPath); statErr != nil {
		return
	}
	var f *os.File
	f, err = os.Open(historyPath)
	if err != nil {
		return
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(&output)
	return
}
