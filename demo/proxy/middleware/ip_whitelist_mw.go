package middleware

import "net"

//IP白名单
func IpWhiteListMiddleWare() func(c *SliceRouterContext) {
	return func(c *SliceRouterContext) {
		remoteAddr, _, _ := net.SplitHostPort(c.Req.RemoteAddr)

		if remoteAddr == "127.0.0.1" {
			c.Next()
			return
		}

		c.Abort()
		c.Rw.Write([]byte("ip_whitelist auth invalid"))
	}
}
