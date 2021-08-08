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
	rb := load_balance.LoadBanlanceFactory(load_balance.LbWeightRoundRobin)
	_ = rb.Add("http://127.0.0.1:2003", "50")

	proxy := proxy3.NewLoadBalanceReverseProxy(&middleware.SliceRouterContext{}, rb)

	log.Println("Starting httpserver at " + addr)
	log.Fatal(http.ListenAndServe(addr, proxy))

}
