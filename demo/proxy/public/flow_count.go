package public

import (
	"fmt"
	"sync/atomic"
	"time"
)

type FlowCountService struct {
	AppId       string        //应用id
	Interval    time.Duration //采集频率
	TotalCount  int64         //当前请求总次数
	QPS         int64         //qps
	Unix        int64         //上次 unix 时间戳
	TickerCount int64         //当前流量
}

func NewFlowCountService(appId string, interval time.Duration) (*FlowCountService, error) {
	reqCounter := &FlowCountService{
		AppId:       appId,
		Interval:    interval,
		QPS:         0,
		Unix:        0,
		TickerCount: 0,
	}

	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()

		ticker := time.NewTicker(interval)
		for true {
			<-ticker.C

			//原子操作
			tickerCount := atomic.LoadInt64(&reqCounter.TickerCount)
			atomic.StoreInt64(&reqCounter.TickerCount, 0)
			nowUnix := time.Now().Unix()

			if reqCounter.Unix == 0 {
				reqCounter.Unix = time.Now().Unix()
				continue
			}

			if nowUnix > reqCounter.Unix {
				reqCounter.QPS = tickerCount / (nowUnix - reqCounter.Unix)
				reqCounter.TotalCount = reqCounter.TotalCount + tickerCount
				reqCounter.Unix = time.Now().Unix()
			}
		}
	}()

	return reqCounter, nil
}

//原子增加
func (o *FlowCountService) Increase() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()

		atomic.AddInt64(&o.TickerCount, 1)
	}()
}
