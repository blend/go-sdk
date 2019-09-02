package jobkit

import (
	"fmt"

	"github.com/blend/go-sdk/cron"
	"github.com/blend/go-sdk/web"
)

// NewManagementServer returns a new management server that lets you
// trigger jobs or look at job statuses via. a json api.
func NewManagementServer(jm *cron.JobManager, cfg Config, options ...web.Option) *web.App {
	app := web.MustNew(append([]web.Option{web.OptConfig(cfg.Web)}, options...)...)
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
	app.GET("/healthz", func(_ *web.Ctx) web.Result {
		if jm.IsStarted() {
			return web.JSON.OK()
		}
		return web.JSON.InternalError(fmt.Errorf("job manager is stopped or in an inconsistent state"))
	})
	app.GET("/api/jobs", func(_ *web.Ctx) web.Result {
		return web.JSON.Result(jm.Status())
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
		invocation := job.GetInvocationByID(web.StringValue(r.RouteParam("invocation")))
		if invocation == nil {
			return web.JSON.NotFound()
		}
		return web.JSON.Result(invocation)
	})
	return app
}
