/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/blend/go-sdk/certutil"
	"github.com/blend/go-sdk/graceful"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/proxyprotocol"
	"github.com/blend/go-sdk/reverseproxy"
	"github.com/blend/go-sdk/webutil"
)

func main() {
	log, err := logger.New(
		logger.OptConfigFromEnv(),
		logger.OptEnabled(webutil.FlagHTTPRequest),
		logger.OptPath("reverse-proxy"),
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

	var listenerOptions []proxyprotocol.CreateListenerOption

	proxy, _ := reverseproxy.NewProxy()
	proxy.Log = log

	var servers []graceful.Graceful
	for _, upstream := range upstreams {
		log.Infof("upstream: %s", upstream)
		target, err := url.Parse(upstream)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		proxyUpstream := reverseproxy.NewUpstream(target)
		proxyUpstream.Log = log
		if err = proxyUpstream.UseHTTP2(); err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
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
	if len(tlsCert) > 0 && len(tlsKey) > 0 {
		certFileWatcher, err := certutil.NewCertFileWatcher(certutil.KeyPair{CertPath: tlsCert, KeyPath: tlsKey})
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		log.Infof("watching tls cert/key files for changes")
		servers = append(servers, certFileWatcher)

		rootCAs, err := certutil.ExtendSystemCertPool()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		proxyServerTLSConfig := &tls.Config{
			RootCAs:        rootCAs,
			GetCertificate: certFileWatcher.GetCertificate,
		}
		webutil.TLSSecureCipherSuites(proxyServerTLSConfig)
		listenerOptions = append(listenerOptions, proxyprotocol.OptTLSConfig(proxyServerTLSConfig))
	}

	proxyServerListener, err := proxyprotocol.CreateListener("tcp", addr, listenerOptions...)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	proxyServer := &http.Server{
		Handler: webutil.NestMiddleware(proxy.ServeHTTP, webutil.HTTPLogged(log)),
	}
	servers = append(servers,
		webutil.NewGracefulHTTPServer(proxyServer, webutil.OptGracefulHTTPServerListener(proxyServerListener)),
	)

	if upgradeAddr != "" {
		log.Infof("http upgrader listening on: %s", upgradeAddr)
		upgrader := reverseproxy.HTTPRedirect{}
		servers = append(servers, webutil.NewGracefulHTTPServer(&http.Server{
			Addr:    upgradeAddr,
			Handler: webutil.NestMiddleware(upgrader.ServeHTTP, webutil.HTTPLogged(log)),
		}))
	}

	log.Infof("reverse proxy listening on: %s", addr)
	if err := graceful.Shutdown(
		servers...,
	); err != nil {
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
