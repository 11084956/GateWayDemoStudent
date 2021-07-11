package circuit_breaker

import (
	"GateWayDemoStudent/demo/proxy/middleware"
	"GateWayDemoStudent/demo/proxy/proxy"
	"GateWayDemoStudent/demo/proxy/public"
	"log"
	"net/http"
	"net/url"
)

var addr = "127.0.0.1:2002"

//熔断方案
func main() {
	coreFunc := func(c *middleware.SliceRouterContext) http.Handler {
		var (
			rs1 = "http://127.0.0.1/2003/base"
			rs2 = "http://127.0.0.1:2004/base"
		)

		url1, err1 := url.Parse(rs1)
		if err1 != nil {
			log.Println(err1)
		}

		url2, err2 := url.Parse(rs2)
		if err2 != nil {
			log.Println(err2)
		}

		urls := []*url.URL{url1, url2}

		return proxy.NewMultipleHostsReverseProxy(c, urls)
	}

	log.Println("starting httpserver at " + addr)
	public.ConfCricuitBreaker(true)

	sliceRouter := middleware.NewSliceRouter()

	sliceRouter.Group("/").Use(middleware.CircuitMW())

	routerHandler := middleware.NewSliceRouterHandler(coreFunc, sliceRouter)

	log.Fatal(http.ListenAndServe(addr, routerHandler))
}
