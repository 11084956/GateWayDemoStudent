package main

import (
	"GateWayDemoStudent/demo/proxy/middleware"
	"GateWayDemoStudent/proxy/load_balance"
	proxy3 "GateWayDemoStudent/proxy/proxy"
	"log"
	"net/http"
)

var addr = "127.0.0.1:2002"

func main() {
	mConf, err := load_balance.NewLoadBalanceCheckConf(
		"http://%s/base",
		map[string]string{"127.0.0.1:2003": "20", "127.0.0.1:2004": "20"})

	if err != nil {
		panic(err)
	}

	rb := load_balance.LoadBanlanceFactorWithConf(load_balance.LbWeightRoundRobin, mConf)
	proxy := proxy3.NewLoadBalanceReverseProxy(&middleware.SliceRouterContext{}, rb)

	log.Println("Starting httpserver at " + addr)
	log.Fatal(http.ListenAndServe(addr, proxy))
}
