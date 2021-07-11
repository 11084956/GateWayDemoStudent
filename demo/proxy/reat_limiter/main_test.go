package reat_limiter

import (
	"context"
	"golang.org/x/time/rate"
	"log"
	"testing"
	"time"
)

func TestRateLimiter(t *testing.T) {
	l := rate.NewLimiter(1, 5)

	log.Println(l.Limit(), l.Burst())

	for i := 0; i < 100; i++ {
		//阻塞等待, 知道取到下一个token
		log.Println("before wait")
		c, _ := context.WithTimeout(context.Background(), time.Second*2)

		if err := l.Wait(c); err != nil {
			log.Println("limiter wait err:" + err.Error())
		}

		log.Println("after wait")

		//返回需要等待多久才会有新的token, 这样就可以在等待时间执行任务

		r := l.Reserve()
		log.Println("reserve Delay:", r.Delay())

		//判断当前是否可以取到 token

		a := l.Allow()
		log.Println("allow:", a)
	}
}
