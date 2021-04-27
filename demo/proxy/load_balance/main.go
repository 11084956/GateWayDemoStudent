package main

import (
	"net"
	"net/http"
	"strings"
	"time"
)

var (
	addr      = "127.0.0.1:2002"
	transport = &http.Transport{
		DialContext: (net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
)

func NewMultipleHostsReverseProxy()  {

}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")

	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "+" + b
	}

	return a + b
}
