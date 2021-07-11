package circuit_breaker

import (
	"errors"
	"github.com/afex/hystrix-go/hystrix"
	"log"
	"net/http"
	"testing"
	"time"
)

func Test_aaMain(m *testing.T) {
	handler := hystrix.NewStreamHandler()
	handler.Start()

	go http.ListenAndServe(":8074", handler)
	hystrix.ConfigureCommand("aaa", hystrix.CommandConfig{
		Timeout:                1000, //单次请求超时时间
		MaxConcurrentRequests:  1,    //最大并发量
		SleepWindow:            5000, //熔断多久后去尝试服务是否可用
		RequestVolumeThreshold: 1,    //验证熔断的请求数量, 10秒内采样
		ErrorPercentThreshold:  1,    //验证熔断的 错误百分比
	})

	for i := 0; i < 10000; i++ {
		//异步调用使用 hystrix.Go
		err := hystrix.Do("aaa", func() error {
			//test case 1并发测试
			if i == 0 {
				return errors.New("server error")
			}

			//test case 2 超时测试
			log.Println("do services")
			return nil
		}, nil)

		if err != nil {
			log.Println("hystrix err:" + err.Error())

			time.Sleep(1 * time.Second)

			log.Println("sleep 1 second")
		}
	}

	time.Sleep(100 * time.Second)
}
