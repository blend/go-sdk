package migration

import (
	"context"
	"fmt"

	"github.com/blend/go-sdk/db"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/logger"
)

// New returns a new suite of groups.
func New(groups ...GroupedActions) *Suite {
	return &Suite{
		Groups: groups,
	}
}

// Suite is a migration suite.
type Suite struct {
	Log    logger.Log
	Groups []GroupedActions

	Applied int
	Skipped int
	Failed  int
	Total   int
}

// Apply applies the suite.
func (s *Suite) Apply(ctx context.Context, c *db.Connection) (err error) {
	defer s.WriteStats(ctx)
	defer func() {
		if r := recover(); r != nil {
			err = ex.New(r)
		}
	}()

	for _, group := range s.Groups {
		if err = group.Action(WithSuite(ctx, s), c); err != nil {
			return
		}
	}
	return
}

// Applyf writes an applied step message.
func (s *Suite) Applyf(ctx context.Context, format string, args ...interface{}) {
	s.Applied = s.Applied + 1
	s.Total = s.Total + 1
	s.Write(ctx, StatApplied, fmt.Sprintf(format, args...))
}

// Skipf skips a given step.
func (s *Suite) Skipf(ctx context.Context, format string, args ...interface{}) {
	s.Skipped = s.Skipped + 1
	s.Total = s.Total + 1
	s.Write(ctx, StatSkipped, fmt.Sprintf(format, args...))
}

// Errorf writes an error for a given step.
func (s *Suite) Errorf(ctx context.Context, format string, args ...interface{}) {
	s.Failed = s.Failed + 1
	s.Total = s.Total + 1
	s.Write(ctx, StatFailed, fmt.Sprintf(format, args...))
}

// Error
func (s *Suite) Error(ctx context.Context, err error) error {
	s.Failed = s.Failed + 1
	s.Total = s.Total + 1
	s.Write(ctx, StatFailed, fmt.Sprintf("%v", err))
	return err
}

func (s *Suite) Write(ctx context.Context, result, body string) {
	logger.MaybeTrigger(ctx, s.Log, NewEvent(result, body, GetContextLabels(ctx)...))
}

// WriteStats writes the stats if a logger is configured.
func (s *Suite) WriteStats(ctx context.Context) {
	logger.MaybeTrigger(ctx, s.Log, NewStatsEvent(s.Applied, s.Skipped, s.Failed, s.Total))
}
