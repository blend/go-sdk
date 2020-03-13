package gobrake

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	tdigest "github.com/caio/go-tdigest"
)

const flushPeriod = 15 * time.Second

type RouteInfo struct {
	Environment string
	Method      string
	Route       string
	StatusCode  int
	Start       time.Time
	End         time.Time
}

type routeKey struct {
	Environment string    `json:"environment"`
	Method      string    `json:"method"`
	Route       string    `json:"route"`
	StatusCode  int       `json:"statusCode"`
	Time        time.Time `json:"time"`
}

type routeStat struct {
	mu      sync.Mutex
	Count   int     `json:"count"`
	Sum     float64 `json:"sum"`
	Sumsq   float64 `json:"sumsq"`
	TDigest []byte  `json:"tdigest"`
	td      *tdigest.TDigest
}

func (s *routeStat) Add(ms float64) error {
	if s.td == nil {
		td, err := tdigest.New(tdigest.Compression(20))
		if err != nil {
			return err
		}
		s.td = td
	}

	s.Count++
	s.Sum += ms
	s.Sumsq += ms * ms
	return s.td.Add(ms)
}

type routeKeyStat struct {
	routeKey
	*routeStat
}

type routeFilter func(*RouteInfo) *RouteInfo

// routeStats aggregates information about requests and periodically sends
// collected data to Airbrake.
type routeStats struct {
	opt    *NotifierOptions
	apiURL string

	flushTimer *time.Timer
	addWG      *sync.WaitGroup

	filters []routeFilter

	mu sync.Mutex
	m  map[routeKey]*routeStat
}

func newRouteStats(opt *NotifierOptions) *routeStats {
	return &routeStats{
		opt: opt,
		apiURL: fmt.Sprintf("%s/api/v5/projects/%d/routes-stats",
			opt.Host, opt.ProjectId),
	}
}

func (s *routeStats) init() {
	if s.flushTimer == nil {
		s.flushTimer = time.AfterFunc(flushPeriod, s.Flush)
		s.addWG = new(sync.WaitGroup)
		s.m = make(map[routeKey]*routeStat)
	}
}

// Flush sends to Airbrake route stats.
func (s *routeStats) Flush() {
	s.mu.Lock()

	s.flushTimer = nil
	addWG := s.addWG
	s.addWG = nil
	m := s.m
	s.m = nil

	s.mu.Unlock()

	if m == nil {
		return
	}

	addWG.Wait()
	err := s.send(m)
	if err != nil {
		logger.Printf("routeStats.send failed: %s", err)
	}
}

type routesStatsJSONRequest struct {
	Routes []routeKeyStat `json:"routes"`
}

func (s *routeStats) send(m map[routeKey]*routeStat) error {
	var routes []routeKeyStat
	for k, v := range m {
		err := v.td.Compress()
		if err != nil {
			return err
		}

		b, err := v.td.AsBytes()
		if err != nil {
			return err
		}
		v.TDigest = b

		routes = append(routes, routeKeyStat{
			routeKey:  k,
			routeStat: v,
		})
	}

	jsonReq := routesStatsJSONRequest{
		Routes: routes,
	}

	buf := buffers.Get().(*bytes.Buffer)
	defer buffers.Put(buf)

	buf.Reset()
	err := json.NewEncoder(buf).Encode(jsonReq)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", s.apiURL, buf)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+s.opt.ProjectKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.opt.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	buf.Reset()
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	switch resp.StatusCode {
	case http.StatusUnauthorized:
		return errUnauthorized
	}

	err = fmt.Errorf("got unexpected response status=%q", resp.Status)
	return err
}

// Notify adds new route stats.
func (s *routeStats) Notify(req *RouteInfo) error {
	for _, fn := range s.filters {
		req = fn(req)

		if req == nil {
			return nil
		}
	}

	key := routeKey{
		Environment: req.Environment,
		Method:      req.Method,
		Route:       req.Route,
		StatusCode:  req.StatusCode,
		Time:        req.Start.UTC().Truncate(time.Minute),
	}

	s.mu.Lock()
	s.init()
	stat, ok := s.m[key]
	if !ok {
		stat = &routeStat{}
		s.m[key] = stat
	}
	addWG := s.addWG
	s.addWG.Add(1)
	s.mu.Unlock()

	ms := float64(req.End.Sub(req.Start)) / float64(time.Millisecond)

	stat.mu.Lock()
	err := stat.Add(ms)
	addWG.Done()
	stat.mu.Unlock()

	return err
}

// AddFilter adds filter that can change route stat or ignore it by returning nil.
func (s *routeStats) AddFilter(fn func(*RouteInfo) *RouteInfo) {
	s.filters = append(s.filters, fn)
}
