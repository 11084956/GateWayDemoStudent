package main

import (
	"GateWayDemoStudent/demo/proxy/reverse_proxy_https/public"
	"GateWayDemoStudent/demo/proxy/reverse_proxy_https/testdata"
	"fmt"
	"golang.org/x/net/http2"
	"log"
	"net/http"
	"net/url"
	"time"
)

var addr = "example1.com:3002"

func main() {
	rs1 := "https://example1.com:3003"
	url1, err1 := url.Parse(rs1)
	if err1 != nil {
		log.Println(err1)
	}

	urls := []*url.URL{url1}
	proxy := public.NewMultipleHostsReverseProxy(urls)

	log.Println("Starting httpserver at " + addr)

	mux := http.NewServeMux()
	mux.Handle("/", proxy)
	server := &http.Server{
		Addr:         addr,
		WriteTimeout: time.Second * 3,
		Handler:      mux,
	}

	err := http2.ConfigureServer(server, &http2.Server{})
	if err != nil {

		fmt.Println("https err:", err)
		return
	}

	err = server.ListenAndServeTLS(testdata.Path("server.crt"),
		testdata.Path("server.key"))
	if err != nil {

		fmt.Println("tls err:", err)
		return
	}

	err = server.ListenAndServe()
	if err != nil {

		fmt.Println("listenAndServe err:", err)
		return
	}
}
