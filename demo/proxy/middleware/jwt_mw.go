package middleware

import (
	"GateWayDemoStudent/demo/proxy/public"
	"strings"
)

func JwtMiddleWare() func(c *SliceRouterContext) {
	return func(c *SliceRouterContext) {
		token := c.Req.Header.Get("Authorization")

		token = strings.Replace(token, "bearer ", "", -1)

		if _, err := public.Decode(token); err != nil {
			c.Rw.Write([]byte("jwt auth invalid:" + err.Error()))
			c.Abort()

			return
		}

		c.Next()
	}
}
