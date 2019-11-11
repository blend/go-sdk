package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/blend/go-sdk/graceful"
	"github.com/blend/go-sdk/webutil"

	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/proxy"
)

// linker metadata block
// this block must be present
// it is used by goreleaser
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

const (
	// DefaultPort is the default port the proxy listens on.
	DefaultPort = "8888"
)

func main() {
	log := logger.MustNewFromEnv().WithEnabled(logger.HTTPRequest)

	var upstreams Upstreams
	flag.Var(&upstreams, "upstream", "An upstream server to proxy traffic to")

	var tlsCert string
	flag.StringVar(&tlsCert, "tls-cert", "", "The path to the tls certificate file (--tls-key must also be set)")

	var tlsKey string
	flag.StringVar(&tlsKey, "tls-key", "", "The path to the tls key file (--tls-cert must also be set)")

	var bindAddr string
	flag.StringVar(&bindAddr, "listen", ":8080", "The address to listen on.")

	var upstreamHeaders UpstreamHeader
	flag.Var(&upstreamHeaders, "upstream-header", "Upstream heaeders to add for all requests.")

	var logEvents string
	flag.StringVar(&logEvents, "log-events", "", "Logger events to enable or disable. Coalesced with `LOG_EVENTS`")

	flag.Parse()

	if len(upstreams) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	if len(logEvents) > 0 {
		eventSet := logger.NewFlagSetFromValues(strings.Split(logEvents, ",")...)
		log.Flags().CoalesceWith(eventSet)
	}

	reverseProxy := proxy.New().WithLogger(log)

	for _, upstream := range upstreams {
		log.SyncInfof("upstream: %s", upstream)
		target, err := url.Parse(upstream)
		if err != nil {
			log.SyncFatalExit(err)
		}

		proxyUpstream := proxy.NewUpstream(target).
			WithLogger(log)

		reverseProxy.WithUpstream(proxyUpstream)
	}

	for _, header := range upstreamHeaders {
		pieces := strings.SplitN(header, "=", 2)
		if len(pieces) < 2 {
			log.SyncFatalExit(fmt.Errorf("invalid header; must be in the form key=value"))
		}
		log.SyncInfof("proxy using upstream header: %s=%s", pieces[0], pieces[1])
		reverseProxy.WithUpstreamHeader(pieces[0], pieces[1])
	}

	if len(tlsCert) > 0 && len(tlsKey) == 0 {
		log.SyncFatalExit(fmt.Errorf("`--tls-key` is unset, cannot continue"))
	}
	if len(tlsCert) == 0 && len(tlsKey) > 0 {
		log.SyncFatalExit(fmt.Errorf("`--tls-key` is unset, cannot continue"))
	}

	server := &http.Server{}
	server.Handler = reverseProxy

	gs := webutil.NewGracefulServer(server)

	listener, err := net.Listen("tcp", bindAddr)
	if err != nil {
		log.SyncFatalExit(err)
	}

	if len(tlsCert) > 0 && len(tlsKey) > 0 {
		log.SyncInfof("proxy using tls cert: %s", tlsCert)
		log.SyncInfof("proxy using tls key: %s", tlsKey)
		cert, err := tls.LoadX509KeyPair(tlsCert, tlsKey)
		if err != nil {
			log.SyncFatalExit(err)
		}
		gs.WithListener(tls.NewListener(listener, &tls.Config{
			Certificates: []tls.Certificate{cert},
		}))
	} else {
		gs.WithListener(listener)
	}

	log.SyncInfof("proxy listening: %s", bindAddr)
	if err := graceful.Shutdown(gs); err != nil {
		log.SyncFatalExit(err)
	}
}

// Upstreams is a flag variable for upstreams.
type Upstreams []string

// String returns a string representation of the upstreams.
func (u *Upstreams) String() string {
	if u == nil {
		return "<nil>"
	}
	return strings.Join(*u, ", ")
}

// Set adds a flag value.
func (u *Upstreams) Set(value string) error {
	*u = append(*u, value)
	return nil
}

// UpstreamHeader is a flag variable for upstreams.
type UpstreamHeader []string

// String returns a string representation of the upstreams.
func (u *UpstreamHeader) String() string {
	if u == nil {
		return "<nil>"
	}
	return strings.Join(*u, ", ")
}

// Set adds a flag value.
func (u *UpstreamHeader) Set(value string) error {
	*u = append(*u, value)
	return nil
}
