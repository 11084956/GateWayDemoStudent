package public

import (
	"GateWayDemoStudent/demo/proxy/reverse_proxy_https/testdata"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"golang.org/x/net/http2"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

var transport = &http.Transport{
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext,
	//TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
	TLSClientConfig: func() *tls.Config {
		pool := x509.NewCertPool()
		caCertPath := testdata.Path("ca.crt")
		caCrt, _ := ioutil.ReadFile(caCertPath)
		pool.AppendCertsFromPEM(caCrt)

		return &tls.Config{
			RootCAs: pool,
		}
	}(),
	MaxIdleConns:          100,              //最大空闲连接
	IdleConnTimeout:       90 * time.Second, //空闲超时时间
	TLSHandshakeTimeout:   10 * time.Second, //tls握手超时时间
	ExpectContinueTimeout: 1 * time.Second,  //100-continue 超时时间
}

func NewMultipleHostsReverseProxy(targets []*url.URL) *httputil.ReverseProxy {
	//请求协调者

	director := func(req *http.Request) {
		targetIndex := rand.Intn(len(targets))
		target := targets[targetIndex]
		targetQuery := target.RawQuery

		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)

		req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		}

		if _, ok := req.Header["User-Agent"]; !ok {
			req.Header.Set("User-Agent", "user-agent")
		}

	}

	err := http2.ConfigureTransport(transport)
	if err != nil {
		fmt.Println("http2请求失败:", err)
		return nil
	}

	return &httputil.ReverseProxy{
		Director:  director,
		Transport: transport,
	}
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
