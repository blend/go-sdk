package gobrake

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type QueryInfo struct {
	Environment string
	Method      string
	Route       string
	Query       string
	Func        string
	File        string
	Line        int
	Start       time.Time
	End         time.Time
}

type queryKey struct {
	Environment string    `json:"environment"`
	Method      string    `json:"method"`
	Route       string    `json:"route"`
	Query       string    `json:"query"`
	Func        string    `json:"function"`
	File        string    `json:"file"`
	Line        int       `json:"line"`
	Time        time.Time `json:"time"`
}

type queryKeyStat struct {
	queryKey
	*routeStat
}

type QueryStats struct {
	opt    *NotifierOptions
	apiURL string

	flushTimer *time.Timer
	addWG      *sync.WaitGroup

	mu sync.Mutex
	m  map[queryKey]*routeStat
}

func newQueryStats(opt *NotifierOptions) *QueryStats {
	return &QueryStats{
		opt: opt,
		apiURL: fmt.Sprintf("%s/api/v5/projects/%d/queries-stats",
			opt.Host, opt.ProjectId),
	}
}

func (s *QueryStats) init() {
	if s.flushTimer == nil {
		s.flushTimer = time.AfterFunc(flushPeriod, s.flush)
		s.addWG = new(sync.WaitGroup)
		s.m = make(map[queryKey]*routeStat)
	}
}

func (s *QueryStats) flush() {
	s.mu.Lock()

	s.flushTimer = nil
	addWG := s.addWG
	s.addWG = nil
	m := s.m
	s.m = nil

	s.mu.Unlock()

	addWG.Wait()
	err := s.send(m)
	if err != nil {
		logger.Printf("queryStats.send failed: %s", err)
	}
}

type queriesStatsJSONRequest struct {
	Queries []queryKeyStat `json:"queries"`
}

func (s *QueryStats) send(m map[queryKey]*routeStat) error {
	var queries []queryKeyStat
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

		queries = append(queries, queryKeyStat{
			queryKey:  k,
			routeStat: v,
		})
	}

	jsonReq := queriesStatsJSONRequest{
		Queries: queries,
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

func (s *QueryStats) Notify(q *QueryInfo) error {
	key := queryKey{
		Environment: q.Environment,
		Method:      q.Method,
		Route:       q.Route,
		Query:       q.Query,
		Func:        q.Func,
		File:        q.File,
		Line:        q.Line,
		Time:        q.Start.UTC().Truncate(time.Minute),
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

	ms := float64(q.End.Sub(q.Start)) / float64(time.Millisecond)

	stat.mu.Lock()
	err := stat.Add(ms)
	addWG.Done()
	stat.mu.Unlock()

	return err
}
