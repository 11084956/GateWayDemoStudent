package load_balance

import (
	"errors"
	"fmt"
	"strings"
)

//轮询负载均衡
type RoundRobinBalance struct {
	curIndex int
	rss      []string
	//观察主题
	conf LoadBalanceConf
}

func (r *RoundRobinBalance) Add(params ...string) error {
	if len(params) == 0 {
		return errors.New("param len 1 at least")
	}

	add := params[0]
	r.rss = append(r.rss, add)

	return nil
}

func (r *RoundRobinBalance) Next() string {
	if len(r.rss) == 0 {
		return ""
	}

	lens := len(r.rss) //5
	if r.curIndex >= lens {
		r.curIndex = 0
	}

	curAddr := r.rss[r.curIndex]
	r.curIndex = (r.curIndex + 1) % lens

	return curAddr
}

func (r *RoundRobinBalance) Get(key string) (string, error) {
	return r.Next(), nil
}

func (r *RoundRobinBalance) SetConf(conf LoadBalanceConf) {
	r.conf = conf
}

func (r *RoundRobinBalance) Update() {
	if conf, ok := r.conf.(*LoadBalanceZkConf); ok {
		confList := conf.GetConf()
		fmt.Println("Update get conf:", confList)
		r.rss = []string{}
		for _, ip := range confList {
			_ = r.Add(strings.Split(ip, ",")...)
		}
	}

	if conf, ok := r.conf.(*LoadBalanceCheckConf); ok {
		confList := conf.GetConf()
		fmt.Println("Update get conf:", confList)
		r.rss = nil
		for _, ip := range confList {
			_ = r.Add(strings.Split(ip, ",")...)
		}
	}
}
