package jobkit

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/blend/go-sdk/cron"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/selector"
	"github.com/blend/go-sdk/uuid"
	"github.com/blend/go-sdk/web"
	"github.com/blend/go-sdk/webutil"

	"github.com/blend/go-sdk/jobkit/static"
	"github.com/blend/go-sdk/jobkit/views"
)

// NewServer returns a new management server that lets you
// trigger jobs or look at job statuses via. a json api.
func NewServer(jm *cron.JobManager, cfg Config, options ...web.Option) *web.App {
	options = append([]web.Option{web.OptConfig(cfg.Web)}, options...)
	app := web.MustNew(options...)

	viewPaths := []string{
		"_views/header.html",
		"_views/footer.html",
		"_views/index.html",
		"_views/invocation.html",
	}
	if cfg.UseViewFilesOrDefault() {
		app.Views.LiveReload = true
		app.Views.AddPaths(viewPaths...)
	} else {
		for _, viewPath := range viewPaths {
			vf, err := views.GetBinaryAsset(viewPath)
			if err != nil {
				panic(err)
			}
			contents, err := vf.Contents()
			if err != nil {
				panic(err)
			}
			app.Views.AddLiterals(string(contents))
		}
	}

	addConfigState := func(action web.Action) web.Action {
		return func(r *web.Ctx) web.Result {
			r.State.Set("Config", cfg)
			return action(r)
		}
	}
	app.DefaultMiddleware = append(app.DefaultMiddleware, addConfigState)

	app.PanicAction = func(r *web.Ctx, err interface{}) web.Result {
		return r.Views.InternalError(ex.New(err))
	}
	app.GET("/", func(r *web.Ctx) web.Result {
		return r.Views.View("index", jm.Status())
	})
	app.GET("/search", func(r *web.Ctx) web.Result {
		selectorParam := web.StringValue(r.QueryValue("selector"))
		if selectorParam == "" {
			return web.RedirectWithMethod("GET", "/")
		}
		sel, err := selector.Parse(selectorParam)
		if err != nil {
			return r.Views.BadRequest(err)
		}

		status := jm.Status()
		status.Jobs = cron.FilterJobSchedulers(status.Jobs, func(js *cron.JobScheduler) bool {
			return sel.Matches(js.Labels())
		})
		r.State.Set("selector", sel.String())
		return r.Views.View("index", status)
	})

	app.GET("/static/*filepath", func(r *web.Ctx) web.Result {
		path, err := r.RouteParam("filepath")
		if err != nil {
			web.Text.NotFound()
		}
		path = filepath.Join("_static", path)
		file, err := static.GetBinaryAsset(path)
		if err == os.ErrNotExist {
			return web.Text.NotFound()
		}
		contents, err := file.Contents()
		if err != nil {
			return web.Text.InternalError(err)
		}
		http.ServeContent(r.Response, r.Request, path, time.Unix(file.ModTime, 0), bytes.NewReader(contents))
		return nil
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
	app.GET("/api/job/:jobName", func(r *web.Ctx) web.Result {
		jobName, err := r.RouteParam("jobName")
		if err != nil {
			return web.JSON.BadRequest(err)
		}
		job, err := jm.Job(jobName)
		if err := jm.RunJob(jobName); err != nil {
			return web.JSON.BadRequest(err)
		}
		return web.JSON.Result(job)
	})
	app.GET("/api/job.stats", func(r *web.Ctx) web.Result {
		output := make(map[string]cron.JobStats)
		for _, job := range jm.Jobs {
			output[job.Name()] = job.Stats()
		}
		return web.JSON.Result(output)
	})
	app.GET("/api/job.stats/:jobName", func(r *web.Ctx) web.Result {
		jobName, err := r.RouteParam("jobName")
		if err != nil {
			return web.JSON.BadRequest(err)
		}
		job, err := jm.Job(jobName)
		if err := jm.RunJob(jobName); err != nil {
			return web.JSON.BadRequest(err)
		}
		return web.JSON.Result(job.Stats())
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

	var getJobInvocation = func(r *web.Ctx, resultProvider web.ResultProvider) (*cron.JobInvocation, web.Result) {
		job, err := jm.Job(web.StringValue(r.RouteParam("jobName")))
		if err != nil {
			return nil, resultProvider.BadRequest(err)
		}
		invocationID, err := r.RouteParam("invocation")
		if err != nil {
			return nil, resultProvider.BadRequest(err)
		}

		var invocation *cron.JobInvocation
		if invocationID == "current" && len(job.Current) > 0 {
			for _, invocation = range job.Current {
				break
			}
		} else if current, ok := job.Current[invocationID]; ok {
			invocation = current
		} else {
			invocation = job.GetInvocationByID(invocationID)
		}

		if invocation == nil {
			return nil, resultProvider.NotFound()
		}
		return invocation, nil
	}

	app.GET("/job.invocation/:jobName/:invocation", func(r *web.Ctx) web.Result {
		invocation, result := getJobInvocation(r, r.Views)
		if result != nil {
			return result
		}
		return r.Views.View("invocation", invocation)
	})
	app.GET("/api/job.invocation/:jobName/:invocation", func(r *web.Ctx) web.Result {
		invocation, result := getJobInvocation(r, web.JSON)
		if result != nil {
			return result
		}
		return web.JSON.Result(invocation)
	})
	app.GET("/job.invocation.output/:jobName/:invocation", func(r *web.Ctx) web.Result {
		invocation, result := getJobInvocation(r, web.Text)
		if result != nil {
			return result
		}
		return web.Text.Result(invocation.Output.String())
	})
	app.GET("/api/job.invocation.output/:jobName/:invocation", func(r *web.Ctx) web.Result {
		invocation, result := getJobInvocation(r, web.JSON)
		if result != nil {
			return result
		}
		lines := append(invocation.Output.Lines)
		if !invocation.Output.Working.Timestamp.IsZero() {
			lines = append(lines, invocation.Output.Working)
		}
		if afterNanos, _ := web.Int64Value(r.QueryValue("afterNanos")); afterNanos > 0 {
			afterTS := time.Unix(0, afterNanos)
			lines = cron.FilterOutputLines(lines, func(l cron.OutputLine) bool {
				return l.Timestamp.After(afterTS)
			})
		}
		return web.JSON.Result(map[string]interface{}{
			"serverTimeNanos": time.Now().UTC().UnixNano(),
			"lines":           lines,
		})
	})

	app.GET("/api/job.invocation.output.stream/:jobName/:invocation", func(r *web.Ctx) web.Result {
		invocation, result := getJobInvocation(r, web.JSON)
		if result != nil {
			return result
		}
		if !jm.IsJobInvocationRunning(invocation.JobName, invocation.ID) {
			return nil
		}

		r.Response.Header().Set(webutil.HeaderContentType, "text/event-stream")
		r.Response.Header().Set(webutil.HeaderVary, "Content-Type")
		r.Response.WriteHeader(http.StatusOK)

		io.WriteString(r.Response, "event: ping\n\n")

		listenerID := uuid.V4().String()

		shouldClose := make(chan struct{})
		invocation.OutputListeners.Add(listenerID, func(l cron.OutputLine) {
			io.WriteString(r.Response, "data: ")
			if _, err := r.Response.Write([]byte(string(l.Data) + "\n\n")); err != nil {
				logger.MaybeError(app.Log, err)
				if shouldClose != nil {
					close(shouldClose)
					shouldClose = nil
				}
			}
			r.Response.Flush()
		})
		defer func() { invocation.OutputListeners.Remove(listenerID) }()

		alarm := time.Tick(500 * time.Millisecond)
		for {
			select {
			case <-shouldClose:
				return nil
			case <-alarm:
				if !jm.IsJobInvocationRunning(invocation.JobName, invocation.ID) {
					return nil
				}
				io.WriteString(r.Response, "event: ping\n\n")
				r.Response.Flush()
			}
		}
	})
	return app
}
