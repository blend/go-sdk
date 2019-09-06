package jobkit

import (
	"fmt"
	"time"

	"github.com/blend/go-sdk/cron"
	"github.com/blend/go-sdk/stringutil"
	"github.com/blend/go-sdk/web"
)

// NewManagementServer returns a new management server that lets you
// trigger jobs or look at job statuses via. a json api.
func NewManagementServer(jm *cron.JobManager, cfg Config, options ...web.Option) *web.App {
	app := web.MustNew(append([]web.Option{web.OptConfig(cfg.Web)}, options...)...)
	app.DefaultMiddleware = nil
	app.Views.AddLiterals(
		headerTemplate,
		footerTemplate,
		indexTemplate,
		invocationTemplate,
	)
	app.GET("/", func(r *web.Ctx) web.Result {
		return r.Views.View("index", jm.Status())
	})
	app.GET("/status.json", func(r *web.Ctx) web.Result {
		return web.JSON.Result(jm.Status())
	})
	app.GET("/api/jobs", func(_ *web.Ctx) web.Result {
		return web.JSON.Result(jm.Status().Jobs)
	})
	app.GET("/api/jobs.running", func(_ *web.Ctx) web.Result {
		return web.JSON.Result(jm.Status().Running)
	})
	app.POST("/pause", func(_ *web.Ctx) web.Result {
		jm.Pause()
		return web.RedirectWithMethod("GET", "/")
	})
	app.POST("/resume", func(_ *web.Ctx) web.Result {
		jm.Resume()
		return web.RedirectWithMethod("GET", "/")
	})
	app.POST("/api/pause", func(_ *web.Ctx) web.Result {
		jm.Pause()
		return web.JSON.OK()
	})
	app.POST("/api/resume", func(_ *web.Ctx) web.Result {
		jm.Resume()
		return web.JSON.OK()
	})
	app.GET("/api/job.status/:jobName", func(r *web.Ctx) web.Result {
		jobName, err := r.RouteParam("jobName")
		if err != nil {
			return web.JSON.BadRequest(err)
		}
		status, err := jm.Job(jobName)
		if err := jm.RunJob(jobName); err != nil {
			return web.JSON.BadRequest(err)
		}
		return web.JSON.Result(status)
	})
	app.POST("/job.run/:jobName", func(r *web.Ctx) web.Result {
		jobName, err := r.RouteParam("jobName")
		if err != nil {
			return r.Views.BadRequest(err)
		}
		if err := jm.RunJob(jobName); err != nil {
			return r.Views.BadRequest(err)
		}
		return web.RedirectWithMethod("GET", "/")
	})
	app.POST("/api/job.run/:jobName", func(r *web.Ctx) web.Result {
		jobName, err := r.RouteParam("jobName")
		if err != nil {
			return web.JSON.BadRequest(err)
		}
		if err := jm.RunJob(jobName); err != nil {
			return web.JSON.BadRequest(err)
		}
		return web.JSON.OK()
	})
	app.POST("/api/job.cancel/:jobName", func(r *web.Ctx) web.Result {
		jobName, err := r.RouteParam("jobName")
		if err != nil {
			return web.JSON.BadRequest(err)
		}
		if err := jm.CancelJob(jobName); err != nil {
			return web.JSON.BadRequest(err)
		}
		return web.JSON.OK()
	})
	app.POST("/job.cancel/:jobName", func(r *web.Ctx) web.Result {
		jobName, err := r.RouteParam("jobName")
		if err != nil {
			return r.Views.BadRequest(err)
		}
		if err := jm.CancelJob(jobName); err != nil {
			return r.Views.BadRequest(err)
		}
		return web.RedirectWithMethod("GET", "/")
	})
	app.POST("/api/job.disable/:jobName", func(r *web.Ctx) web.Result {
		jobName, err := r.RouteParam("jobName")
		if err != nil {
			return web.JSON.BadRequest(err)
		}
		if err := jm.DisableJobs(jobName); err != nil {
			return web.JSON.BadRequest(err)
		}
		return web.JSON.Result(fmt.Sprintf("%s disabled", jobName))
	})
	app.POST("/job.disable/:jobName", func(r *web.Ctx) web.Result {
		jobName, err := r.RouteParam("jobName")
		if err != nil {
			return r.Views.BadRequest(err)
		}
		if err := jm.DisableJobs(jobName); err != nil {
			return r.Views.BadRequest(err)
		}
		return web.RedirectWithMethod("GET", "/")
	})
	app.POST("/api/job.enable/:jobName", func(r *web.Ctx) web.Result {
		jobName, err := r.RouteParam("jobName")
		if err != nil {
			return web.JSON.BadRequest(err)
		}
		if err := jm.EnableJobs(jobName); err != nil {
			return web.JSON.BadRequest(err)
		}
		return web.JSON.Result(fmt.Sprintf("%s enabled", jobName))
	})
	app.POST("/job.enable/:jobName", func(r *web.Ctx) web.Result {
		jobName, err := r.RouteParam("jobName")
		if err != nil {
			return r.Views.BadRequest(err)
		}
		if err := jm.EnableJobs(jobName); err != nil {
			return r.Views.BadRequest(err)
		}
		return web.RedirectWithMethod("GET", "/")
	})
	app.GET("/job.invocation/:jobName/:invocation", func(r *web.Ctx) web.Result {
		job, err := jm.Job(web.StringValue(r.RouteParam("jobName")))
		if err != nil {
			return r.Views.BadRequest(err)
		}
		invocationID, err := r.RouteParam("invocation")
		if err != nil {
			return r.Views.BadRequest(err)
		}
		if job.Current != nil && job.Current.ID == invocationID {
			return r.Views.View("invocation", job.Current)
		}
		invocation := job.GetInvocationByID(invocationID)
		if invocation == nil {
			return r.Views.NotFound()
		}
		return r.Views.View("invocation", invocation)
	})
	app.GET("/api/job.invocation/:jobName/:invocation", func(r *web.Ctx) web.Result {
		job, err := jm.Job(web.StringValue(r.RouteParam("jobName")))
		if err != nil {
			return web.JSON.BadRequest(err)
		}
		invocationID, err := r.RouteParam("invocation")
		if err != nil {
			return web.JSON.BadRequest(err)
		}
		if job.Current != nil && job.Current.ID == invocationID {
			return web.JSON.Result(job.Current)
		}
		invocation := job.GetInvocationByID(invocationID)
		if invocation == nil {
			return web.JSON.NotFound()
		}
		return web.JSON.Result(invocation)
	})
	app.GET("/job.invocation.output/:jobName/:invocation", func(r *web.Ctx) web.Result {
		job, err := jm.Job(web.StringValue(r.RouteParam("jobName")))
		if err != nil {
			return web.Text.BadRequest(err)
		}
		invocationID, err := r.RouteParam("invocation")
		if err != nil {
			return web.Text.BadRequest(err)
		}
		var invocation *cron.JobInvocation
		if job.Current != nil && job.Current.ID == invocationID {
			invocation = job.Current
		} else {
			invocation = job.GetInvocationByID(invocationID)
		}
		if invocation == nil {
			return web.Text.NotFound()
		}
		return web.Text.Result(invocation.Output.String())
	})
	app.GET("/api/job.invocation.output/:jobName/:invocation", func(r *web.Ctx) web.Result {
		job, err := jm.Job(web.StringValue(r.RouteParam("jobName")))
		if err != nil {
			return web.JSON.BadRequest(err)
		}
		invocationID, err := r.RouteParam("invocation")
		if err != nil {
			return web.JSON.BadRequest(err)
		}
		var invocation *cron.JobInvocation
		if job.Current != nil && job.Current.ID == invocationID {
			invocation = job.Current
		} else {
			invocation = job.GetInvocationByID(invocationID)
		}
		if invocation == nil {
			return web.JSON.NotFound()
		}
		lines := append(invocation.Output.Lines)
		if !invocation.Output.Current.Timestamp.IsZero() {
			lines = append(lines, invocation.Output.Current)
		}
		if afterNanos, _ := web.Int64Value(r.QueryValue("afterNanos")); afterNanos > 0 {
			afterTS := time.Unix(0, afterNanos)
			lines = stringutil.FilterLines(lines, func(l stringutil.Line) bool {
				return l.Timestamp.After(afterTS)
			})
		}
		return web.JSON.Result(map[string]interface{}{
			"serverTimeNanos": time.Now().UTC().UnixNano(),
			"lines":           lines,
		})
	})
	return app
}
