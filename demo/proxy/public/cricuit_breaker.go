package public

import (
	"github.com/afex/hystrix-go/hystrix"
	"log"
	"net"
	"net/http"
)

func ConfCricuitBreaker(openStream bool) {
	hystrix.ConfigureCommand("common", hystrix.CommandConfig{
		Timeout:                1000, //单词请求超时时间
		MaxConcurrentRequests:  1,    //最大并发数
		SleepWindow:            5000, //熔断多久后尝试服务是否可用
		RequestVolumeThreshold: 1,    //验证熔断的请求数量
		ErrorPercentThreshold:  1,    //验证熔断的错误百分比
	})

	if openStream {
		handler := hystrix.NewStreamHandler()
		handler.Start()

		go func() {
			err := http.ListenAndServe(net.JoinHostPort("", "2001"), handler)

			log.Fatal(err)
		}()
	}
}
