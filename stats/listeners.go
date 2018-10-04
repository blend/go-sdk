package stats

import (
	"reflect"
	"strconv"

	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/exception"
	"github.com/blend/go-sdk/logger"
)

// AddDefaultTagsFromEnv adds default tags to a collector from environment values.
func AddDefaultTagsFromEnv(stats Collector) {
	if stats == nil {
		return
	}
	stats.AddDefaultTag(TagService, env.Env().String("SERVICE_NAME"))
	stats.AddDefaultTag(TagEnv, env.Env().String("SERVICE_ENV"))
	stats.AddDefaultTag(TagContainer, env.Env().String("HOSTNAME"))
}

// AddWebListeners adds web listeners.
func AddWebListeners(log *logger.Logger, stats Collector) {
	if log == nil || stats == nil || reflect.ValueOf(stats).IsNil() {
		return
	}

	log.Listen(logger.HTTPResponse, ListenerNameStats, logger.NewHTTPResponseEventListener(func(wre *logger.HTTPResponseEvent) {
		var route string
		if len(wre.Route()) > 0 {
			route = Tag(TagRoute, wre.Route())
		} else {
			route = Tag(TagRoute, RouteNotFound)
		}

		method := Tag(TagMethod, wre.Request().Method)
		status := Tag(TagStatus, strconv.Itoa(wre.StatusCode()))
		tags := []string{
			route, method, status,
		}

		stats.Increment(MetricNameHTTPRequest, tags...)
		stats.TimeInMilliseconds(MetricNameHTTPRequestElapsed, wre.Elapsed(), tags...)
	}))
}

// AddQueryListeners adds db listeners.
func AddQueryListeners(log *logger.Logger, stats Collector) {
	if log == nil || stats == nil || reflect.ValueOf(stats).IsNil() {
		return
	}

	log.Listen(logger.Query, ListenerNameStats, logger.NewQueryEventListener(func(qe *logger.QueryEvent) {
		engine := Tag(TagEngine, qe.Engine())
		database := Tag(TagDatabase, qe.Database())

		tags := []string{
			engine, database,
		}
		if len(qe.QueryLabel()) > 0 {
			tags = append(tags, Tag(TagQuery, qe.QueryLabel()))
		}
		if qe.Err() != nil {
			if ex := exception.As(qe.Err()); ex != nil && ex.Class() != nil {
				tags = append(tags, Tag(TagClass, ex.Class().Error()))
			}
			tags = append(tags, TagError)
		}

		stats.Increment(MetricNameDBQuery, tags...)
		stats.TimeInMilliseconds(MetricNameDBQueryElapsed, qe.Elapsed(), tags...)
	}))
}

// AddErrorListeners adds error listeners.
func AddErrorListeners(log *logger.Logger, stats Collector) {
	if log == nil || stats == nil || reflect.ValueOf(stats).IsNil() {
		return
	}

	listener := logger.NewErrorEventListener(func(ee *logger.ErrorEvent) {
		stats.Increment(MetricNameError,
			Tag(TagSeverity, string(ee.Flag())),
			Tag(TagClass, exception.ErrClass(ee.Err())),
		)
	})
	log.Listen(logger.Warning, ListenerNameStats, listener)
	log.Listen(logger.Error, ListenerNameStats, listener)
	log.Listen(logger.Fatal, ListenerNameStats, listener)
}
