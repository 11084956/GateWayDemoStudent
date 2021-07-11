package middleware

import (
	"GateWayDemoStudent/demo/proxy/public"
	"fmt"
)

//统计qps
func FlowCountMiddleWare(counter *public.FlowCountService) func(c *SliceRouterContext) {

	return func(c *SliceRouterContext) {
		counter.Increase()

		fmt.Println("QPS:", counter.QPS)
		fmt.Println("TotalCount:", counter.TotalCount)

		c.Next()
	}
}

//使用redis 统计qps
func RedisFlowCountMiddleWare(counter *public.RedisFlowCountService) func(c *SliceRouterContext) {
	return func(c *SliceRouterContext) {
		counter.Increase()

		fmt.Println("QPS:", counter.QPS)
		fmt.Println("TotalCount:", counter.TotalCount)
		c.Next()
	}
}
