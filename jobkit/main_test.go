package jobkit

import (
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/cron"
	"github.com/blend/go-sdk/uuid"
	"github.com/blend/go-sdk/web"
)

func firstJob(jm *cron.JobManager) *cron.JobScheduler {
	sorted := sortedJobs(jm)
	if len(sorted) > 0 {
		return sorted[0]
	}
	return nil
}

// sortedJobs returns the list of jobs ordered by job name.
func sortedJobs(jm *cron.JobManager) []*cron.JobScheduler {
	var output []*cron.JobScheduler
	for _, js := range jm.Jobs {
		output = append(output, js)
	}
	sort.Sort(cron.JobSchedulersByJobNameAsc(output))
	return output
}

func createTestOutputLine() cron.OutputLine {
	return cron.OutputLine{
		Timestamp: time.Now().UTC(),
		Data:      []byte(uuid.V4()),
	}
}

func createTestCompleteJobInvocation(jobName string, elapsed time.Duration) cron.JobInvocation {
	return cron.JobInvocation{
		ID:       uuid.V4().String(),
		JobName:  jobName,
		Started:  time.Now().UTC(),
		Finished: time.Now().UTC().Add(elapsed),
		State:    cron.JobInvocationStateComplete,
		Elapsed:  elapsed,
		Output: &cron.OutputBuffer{
			Working: createTestOutputLine(),
			Lines: []cron.OutputLine{
				createTestOutputLine(),
				createTestOutputLine(),
				createTestOutputLine(),
				createTestOutputLine(),
			},
		},
	}
}

func createTestFailedJobInvocation(jobName string, elapsed time.Duration) cron.JobInvocation {
	return cron.JobInvocation{
		ID:       uuid.V4().String(),
		JobName:  jobName,
		Started:  time.Now().UTC(),
		Finished: time.Now().UTC().Add(elapsed),
		State:    cron.JobInvocationStateFailed,
		Elapsed:  elapsed,
		Err:      fmt.Errorf("this is only a test: %s", uuid.V4().String()),
		Output: &cron.OutputBuffer{
			Working: createTestOutputLine(),
			Lines: []cron.OutputLine{
				createTestOutputLine(),
				createTestOutputLine(),
			},
		},
	}
}

func createTestManagementServer() (*cron.JobManager, *web.App) {
	test0 := cron.NewJob(cron.OptJobName("test0"))
	test1 := cron.NewJob(cron.OptJobName("test1"))

	jm := cron.New()
	jm.LoadJobs(test0, test1)
	jm.Jobs["test0"].Current = &cron.JobInvocation{
		ID:      uuid.V4().String(),
		JobName: "test0",
		Started: time.Now().UTC(),
		Output: &cron.OutputBuffer{
			Working: createTestOutputLine(),
			Lines: []cron.OutputLine{
				createTestOutputLine(),
				createTestOutputLine(),
				createTestOutputLine(),
				createTestOutputLine(),
			},
		},
	}

	jm.Jobs["test0"].History = []cron.JobInvocation{
		createTestCompleteJobInvocation("test0", 200*time.Millisecond),
		createTestCompleteJobInvocation("test0", 250*time.Millisecond),
		createTestFailedJobInvocation("test0", 5*time.Second),
	}
	jm.Jobs["test1"].History = []cron.JobInvocation{
		createTestCompleteJobInvocation("test1", 200*time.Millisecond),
		createTestCompleteJobInvocation("test1", 250*time.Millisecond),
		createTestCompleteJobInvocation("test1", 300*time.Millisecond),
		createTestCompleteJobInvocation("test1", 350*time.Millisecond),
	}
	return jm, NewServer(jm, Config{
		Web: web.Config{
			Port: 5000,
		},
	})
}

func TestMain(m *testing.M) {
	assert.Main(m)
}
