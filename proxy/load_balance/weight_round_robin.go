package load_balance

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

//加权负载均衡
type WeightRoundRobinBalance struct {
	curIndex int
	rss      []*WeightNode
	rsw      []int

	//观察主体
	conf LoadBalanceConf
}

//忽略了节点故障等
type WeightNode struct {
	addr            string
	weight          int //权重值
	currentWeight   int //节点当前权重
	effectiveWeight int //有效权重
}

func (r *WeightRoundRobinBalance) Add(params ...string) error {
	if len(params) != 2 {
		return errors.New("param len need 2")
	}

	parInt, err := strconv.ParseInt(params[1], 10, 64)
	if err != nil {
		return err
	}

	node := &WeightNode{
		addr:            params[0],
		weight:          int(parInt),
		effectiveWeight: int(parInt),
	}

	r.rss = append(r.rss, node)

	return nil
}

func (r *WeightRoundRobinBalance) Next() string {
	var (
		total = 0
		base  *WeightNode
	)

	for i := 0; i < len(r.rss); i++ {
		w := r.rss[i]
		//step 1.统计所有有效权重之和
		total += w.effectiveWeight

		//step 2.变更节点临时权重为 节点临时权重+节点有效权重
		w.currentWeight += w.effectiveWeight

		//step 3.有效权重默认与权重相同,通讯异常时-1, 通讯成功+1,知道恢复到weight大小
		if w.effectiveWeight < w.weight {
			w.effectiveWeight++
		}

		//step 4.选择最大临时权重节点
		if base == nil || w.currentWeight > base.currentWeight {
			base = w
		}
	}

	if base == nil {
		return ""
	}

	//step 5.变更临时权重为 临时权重-有效权重之和
	base.currentWeight -= total

	return base.addr
}

func (r *WeightRoundRobinBalance) Get() (string, error) {
	return r.Next(), nil
}

func (r *WeightRoundRobinBalance) SetConf(conf LoadBalanceConf) {
	r.conf = conf
}

func (r *WeightRoundRobinBalance) Update() {
	if conf, ok := r.conf.(*LoadBalanceZkConf); ok {
		confList := conf.GetConf()
		fmt.Println("WeightRoundRobinBalance get conf:", confList)
		r.rss = nil
		for _, ip := range confList {
			_ = r.Add(strings.Split(ip, ",")...)
		}
	}

	if conf, ok := r.conf.(*LoadBalanceCheckConf); ok {
		confList := conf.GetConf()
		fmt.Println("WeightRoundRobinBalance get conf:", conf.GetConf())
		r.rss = nil
		for _, ip := range confList {
			_ = r.Add(strings.Split(ip, ",")...)
		}
	}
}
