package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/blend/go-sdk/graceful"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/reverseproxy"
)

// linker metadata block
// this block must be present
// it is used by goreleaser
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	log, err := logger.New(
		logger.OptConfigFromEnv(),
		logger.OptEnabled(logger.HTTPRequest),
		logger.OptEnabled(logger.HTTPResponse),
		logger.OptSubContext("reverse-proxy"),
	)
	if err != nil {
		logger.FatalExit(err)
	}

	var upstreams Upstreams
	flag.Var(&upstreams, "upstream", "An upstream server to proxy traffic to")

	var tlsCert string
	flag.StringVar(&tlsCert, "tls-cert", "", "The path to the tls certificate file (--tls-key must also be set)")

	var tlsKey string
	flag.StringVar(&tlsKey, "tls-key", "", "The path to the tls key file (--tls-cert must also be set)")

	var addr string
	flag.StringVar(&addr, "addr", reverseproxy.DefaultAddr, "The address to listen on.")

	var upgradeAddr string
	flag.StringVar(&upgradeAddr, "upgrade-addr", "", "The upgrade address to listen on.")

	var useProxyProtocol bool
	flag.BoolVar(&useProxyProtocol, "proxyProtocol", false, "If we should decode proxy protocol.")

	var upstreamHeaders UpstreamHeader
	flag.Var(&upstreamHeaders, "upstream-header", "Upstream heaeders to add for all requests.")

	flag.Parse()

	if len(upstreams) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	proxy := reverseproxy.NewProxy()
	for _, upstream := range upstreams {
		log.Infof("upstream: %s", upstream)
		target, err := url.Parse(upstream)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		proxyUpstream := reverseproxy.NewUpstream(target)
		proxyUpstream.Log = log
		proxy.Upstreams = append(proxy.Upstreams, proxyUpstream)
	}

	for _, header := range upstreamHeaders {
		pieces := strings.SplitN(header, "=", 2)
		if len(pieces) < 2 {
			log.Fatal(fmt.Errorf("invalid header; must be in the form key=value"))
			os.Exit(1)
		}
		log.Infof("proxy using upstream header: %s=%s", pieces[0], pieces[1])
		proxy.Headers.Add(pieces[0], pieces[1])
	}

	if len(tlsCert) > 0 && len(tlsKey) == 0 {
		log.Fatal(fmt.Errorf("`--tls-key` is unset, cannot continue"))
		os.Exit(1)
	}
	if len(tlsCert) == 0 && len(tlsKey) > 0 {
		log.Fatal(fmt.Errorf("`--tls-key` is unset, cannot continue"))
		os.Exit(1)
	}

	serverOptions := []reverseproxy.ServerOption{
		reverseproxy.OptServerLog(log),
		reverseproxy.OptServerBindAddr(addr),
		reverseproxy.OptServerUseProxyProtocol(useProxyProtocol),
	}

	if upgradeAddr != "" {
		log.Infof("reverse proxy upgrader listening: %s", upgradeAddr)
		serverOptions = append(serverOptions, reverseproxy.OptServerUpgrade(upgradeAddr))
	}

	if len(tlsCert) > 0 && len(tlsKey) > 0 {
		log.Info("reverse proxy using tls")
		log.Infof("reverse proxy watching tls paths: %s, %s", tlsCert, tlsKey)
		serverOptions = append(serverOptions, reverseproxy.OptServerTLSCertWatcher(tlsCert, tlsKey))
	}

	server, err := reverseproxy.NewServer(proxy, serverOptions...)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	log.Infof("reverse proxy listening: %s", addr)
	if err := graceful.Shutdown(server); err != nil {
		log.Fatal(err)
		os.Exit(1)
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
